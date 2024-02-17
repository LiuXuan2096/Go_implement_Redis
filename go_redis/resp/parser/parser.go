package parser

import (
	"bufio"
	"errors"
	"go_redis/interface/resp"
	"go_redis/lib/logger"
	"go_redis/resp/reply"
	"io"
	"runtime/debug"
	"strconv"
	"strings"
)

/**
 * Redis Serialization Protocol(RESP)规定了Redis客户端和服务端的通信协议，详见README.md
 * 因为Redis客户端和服务端向对方发送的消息都是同样的数据格式，所以代码注释中提到的“回复”
 * 并不仅仅是指 服务端向客户端发送的消息，也可以指客户端向服务端发送的消息，注意理解上不要有歧义。
 *
 * RESP中规定的五种类型的消息格式，其中 “正常回复” “错误回复” “整数” 这三种类型的消息
 * 由 parseSingleLineReply()方法解析
 * “字符串”类型的消息的消息头由parseBulkHeader()方法解析
 * “数组”类型的消息的消息头由parseMultiBulkHeader()方法解析
 * 这两种消息的消息体由readBody()方法解析
 */

// Payload 是存储Redis服务端和Redis客户端相互之间发送信息的数据结构
// 因为Redis客户端和服务端向对方发送的消息的数据格式是一样的，所以可以
// 使用同一个数据结构来表示
type Payload struct {
	Data resp.Reply
	Err  error
}

// ParseStream 从io.Reader中读取数据，底层的TCP Socket返回给上层的
// 就是一个io.Reader，所以我们从io.Reader中读取数据就可以，并将数据解析
// 成 Payload 格式，通过Channel返回给调用方
func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	// 这里我们希望Redis解析客户端发来的消息的工作 和 Redis自身的业务逻辑并发执行
	// 所以开启个新的协程
	go parse0(reader, ch)
	return ch
}

type readState struct {
	// 值为true表示当前解析的不是消息的首行数据
	// 值为false表示当前解析的是消息的首行数据
	readingMultiline bool
	// 表示当前消息应该解析出的参数的数量
	expectedArgsCount int
	// 表示当前解析的消息的类型
	msgType byte
	// 存放解析过程中，解析出的参数 放在这个切片中
	args [][]byte
	// 表示当前解析的语句块的字节数
	bulkLen int64
}

func parse0(reader io.Reader, ch chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	bufReader := bufio.NewReader(reader)
	var state readState
	var err error
	var msg []byte
	for {
		// 读取一行数据
		var ioErr bool
		msg, ioErr, err = readLine(bufReader, &state)
		if err != nil {
			if ioErr {
				// 遇到IO错误，停止读取
				ch <- &Payload{
					Err: err,
				}
				close(ch)
				return
			}
			// 协议错误，重置readState的状态
			ch <- &Payload{
				Err: err,
			}
			state = readState{}
			continue
		}

		// 解析一行数据
		if !state.readingMultiline {
			// readingMultiline这个字段为 false 说明当前消息的是消息的第一行数据
			// 这个字段为true，说明该消息的第一行数据以解析完成，state变量已初始化
			// 接收到新的请求
			if msg[0] == '*' {
				// 消息是以 * 开头的消息
				// 说明要解析的消息含有多个语句块
				// 即消息类型的数组
				err = parseMultiBulkHeader(msg, &state)
				if err != nil {
					// 解析发生错误说明客户端发来的消息格式有问题
					ch <- &Payload{
						Err: errors.New("protocol error: " + string(msg)),
					}
					// 将state置为空
					state = readState{}
					continue
				}
				if state.expectedArgsCount == 0 {
					// 说明客户端发来的是空消息
					ch <- &Payload{
						Data: &reply.EmptyMultiBulkReply{},
					}
					// 将state置为空
					state = readState{}
					continue
				}
			} else if msg[0] == '$' {
				// 消息以 $ 开头，说明要解析的消息只含有一个语句块
				// 即消息类型是 字符串消息
				err = parseBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: errors.New("protocol error: " + string(msg)),
					}
					state = readState{}
					continue
				}
				if state.bulkLen == -1 {
					// null 消息
					ch <- &Payload{
						Data: &reply.NullBulkReply{},
					}
					state = readState{}
					continue
				}
			} else {
				// 能进入这个代码块 说明消息类型是 “正常回复” “错误回复” “整数”
				// 这三种类型的消息，详见 READEME.md
				result, err := parseSingleLineReply(msg)
				ch <- &Payload{
					Data: result,
					Err:  err,
				}
				state = readState{} // 解析完毕后将state变量置为空
				continue
			}
		} else {
			// 执行到这个代码块，说明当前解析的不是消息的首行数据
			// 是首行数据之后的数据
			err = readBody(msg, &state)
			if err != nil {
				ch <- &Payload{
					Err: errors.New("protocol error: " + string(msg)),
				}
				state = readState{}
				continue
			}
			if state.finished() {
				var result resp.Reply
				if state.msgType == '*' {
					result = reply.MakeMultiBulkReply(state.args)
				} else if state.msgType == '$' {
					result = reply.MakeBulkReply(state.args[0])
				}
				ch <- &Payload{
					Data: result,
					Err:  err,
				}
				state = readState{}
			}
		}
	}
}

// finished 判断当前消息是否解析完成
func (s *readState) finished() bool {
	return s.expectedArgsCount > 0 && len(s.args) == s.expectedArgsCount
}

// readLine 从底层的TCP Socket返回的bufio.Reader读取数据，
func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	var msg []byte
	var err error
	if state.bulkLen == 0 {
		// 解析只有含有一个‘\r\n'的消息
		msg, err = bufReader.ReadBytes('\n')
		if err != nil {
			// 返回值中的bool变量值为true，表示解析过程中发生了IO错误
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
	} else {
		// 解析含有多个'\r\n'的消息
		msg = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(bufReader, msg)
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 ||
			msg[len(msg)-2] != '\r' ||
			msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, false, nil
}

// parseMultiBulkHeader 解析数组消息的头部信息。
// 以客户端向Redis服务端发送“Set key value”消息为例，
// 按照RESP协议，该消息实际传输时的文本为：“*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n”
// 其中“*3\r\n”表示本次发送的消息是个数组，该数组有3个成员
// 这个方法的作用就是将 消息头部表示数组成员个数的“3”提取出来
func parseMultiBulkHeader(msg []byte, state *readState) error {
	var err error
	// 这个变量表示的是消息格式为数组类型时，需要处理的成员的个数
	var expectedLine uint64
	// 将消息头部表示数组成员个数的数字从字符串形式解析成uint64的整型
	expectedLine, err = strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 32)
	if err != nil {
		// 如果解析发生错误，说明客户端传来的消息文本的格式有问题
		return errors.New("protocol error: " + string(msg))
	}
	if expectedLine == 0 {
		state.expectedArgsCount = 0
		return nil
	} else if expectedLine > 0 {
		state.msgType = msg[0]
		state.readingMultiline = true
		state.expectedArgsCount = int(expectedLine)
		state.args = make([][]byte, 0, expectedLine)
		return nil
	} else {
		return errors.New("protocol error: " + string(msg))
	}
}

// parseBulkHeader 这个方法是用来解析以$开头的字符串消息的头部信息
// 比如客户端向服务端发送字符串“shishi” 在传输时实际发送的消息的格式是
// $6\r\nShiShi\r\n  以$ 开头 之后的数字表示要发送的消息的字节数
func parseBulkHeader(msg []byte, state *readState) error {
	var err error
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		// 解析出现错误，说明客户端发送过来的消息格式有问题
		return errors.New("protocol error: " + string(msg))
	}
	if state.bulkLen == -1 {
		// 说明消息内容是null
		return nil
	} else if state.bulkLen > 0 {
		state.msgType = msg[0]
		state.readingMultiline = true
		state.expectedArgsCount = 1
		state.args = make([][]byte, 0, 1)
		return nil
	} else {
		return errors.New("protocol error: " + string(msg))
	}
}

// parseSingleLineReply 解析RESP中“正常回复” “错误回复” “整数” 这三种消息格式
// 正常回复：以 + 开头 以 \r\n 结尾的字符串形式  +ok/r/r
// 错误回复：以 - 开头，以 /r/n 结尾的字符串形式 -Error message\r\n
// 整数 以：开头，以/r/n结尾的字符串形式 :123456\r\n
func parseSingleLineReply(msg []byte) (resp.Reply, error) {
	str := strings.TrimSuffix(string(msg), "\r\n")
	var result resp.Reply
	switch msg[0] {
	case '+':
		result = reply.MakeStatusReply(str[1:])
	case '-':
		result = reply.MakeErrReply(str[1:])
	case ':':
		val, err := strconv.ParseInt(str[1:], 10, 64)
		if err != nil {
			return nil, errors.New("protocol error: " + string(msg))
		}
		result = reply.MakeIntReply(val)
	}
	return result, nil
}

// 解析 $ 开头的字符串消息 和 * 开头的数组消息的消息体。
func readBody(msg []byte, state *readState) error {
	line := msg[0 : len(msg)-2]
	var err error
	if line[0] == '$' {
		// 以$开头的单个字符串消息
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return errors.New("protocol error: " + string(msg))
		}
		if state.bulkLen <= 0 {
			// 遇到了null语句块
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}
	return nil
}
