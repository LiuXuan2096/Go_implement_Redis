package reply

import (
	"bytes"
	"go_redis/interface/resp"
	"strconv"
)

var (
	nullBulkReplyBytes = []byte("$-1")

	// CRLF 是RESP(redis serialization protocol)的行分隔符
	CRLF = "\r\n"
)

/*
 * 回复数组，对应RESP中的数组
 */

// BulkReply存储着一个二进制安全的字符串
type BulkReply struct {
	Arg []byte
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{
		Arg: arg,
	}
}

func (r *BulkReply) ToBytes() []byte {
	if len(r.Arg) == 0 {
		return nullBulkBytes
	}
	return []byte("$" + strconv.Itoa(len(r.Arg)) + CRLF + string(r.Arg) + CRLF)
}

/**
 *回复多行字符串，对应RESP中的多行字符串
 */
type MultiBulkReply struct {
	Args [][]byte
}

func MakeMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{
		Args: args,
	}
}

func (r *MultiBulkReply) ToBytes() []byte {
	argLen := len(r.Args)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(argLen) + CRLF)
	for _, arg := range r.Args {
		if arg == nil {
			buf.WriteString("$-1" + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

/**
 * 回复状态,对应RESP中的正常回复
 */
type StatusReply struct {
	Status string
}

func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{
		Status: status,
	}
}

func (r *StatusReply) ToBytes() []byte {
	return []byte("+" + r.Status + CRLF)
}

/*
 * 回复整数值，对应RESP中的回复整数
 */
type IntReply struct {
	Code int64
}

func MakeIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

func (r *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(r.Code, 10) + CRLF)
}

/*
 * 错误回复，对应RESP协议中的错误回复
 */
type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

// 标准错误回复
type StandardErrReply struct {
	Status string
}

func (r *StandardErrReply) ToBytes() []byte {
	return []byte("-" + r.Status + CRLF)
}

func (r *StandardErrReply) Error() string {
	return r.Status
}

func MakeErrReply(status string) *StandardErrReply {
	return &StandardErrReply{
		Status: status,
	}
}

// 如果向客户端发送的消息的类型是“错误消息”，则返回true
func IsErrorReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}
