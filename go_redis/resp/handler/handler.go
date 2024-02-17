package handler

import (
	"context"
	"go_redis/database"
	databaseface "go_redis/interface/database"
	"go_redis/lib/logger"
	"go_redis/lib/sync/atomic"
	"go_redis/resp/connection"
	"go_redis/resp/parser"
	"go_redis/resp/reply"
	"io"
	"net"
	"strings"
	"sync"
)

var (
	// 收到客户端发送的不符合RESP协议的未知消息时，向客户端发送如下回复
	unknowErrReplyBytes = []byte("-ERR unknown\r\n")
)

type RespHandler struct {
	activeConn sync.Map // 存放同Redis客户端的连接的容器
	db         databaseface.Database
	// 表示当前Redis服务端是否处于关闭或正在关闭中
	// 值为true时拒绝新的客户端连接和新的请求，开始执行关闭Redis服务端的逻辑
	closing atomic.Boolean
}

// MakeHandler 返回一个RespHandler实例
func MakeHandler() *RespHandler {
	var db databaseface.Database
	db = database.NewEchoDatabase()
	return &RespHandler{
		db: db,
	}
}

// closeClient 关闭同某个Redis客户端的连接
func (h *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	h.db.AfterClientClose(client)
	h.activeConn.Delete(client)
}

// Handle 接收和执行客户端发来的Redis指令
func (h *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		// 关闭handler即关闭Redis服务端，同时拒绝新的客户端连接
		_ = conn.Close()
	}

	// 包装客户端连接，并将客户端对象放到容器中
	client := connection.NewConn(conn)
	h.activeConn.Store(client, 1)

	// 开始解析客户端发来的指令消息，并将其写入Channel中
	ch := parser.ParseStream(conn)

	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				// 这些错误说明底层的TCP连接已经关闭
				h.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			// 代码执行到这说明是协议错误，底层的TCP连接没有问题
			errReply := reply.MakeErrReply(payload.Err.Error())
			err := client.Write(errReply.ToBytes()) // 将错误信息返回给客户端
			if err != nil {
				h.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			continue
		}
		if payload.Data == nil {
			logger.Error("empty payload")
			continue
		}
		r, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk reply")
			continue
		}
		result := h.db.Exec(client, r.Args)
		if result != nil {
			_ = client.Write(result.ToBytes())
		} else {
			_ = client.Write(unknowErrReplyBytes)
		}
	}
}

// Close 关闭Handler 即关闭Redis的服务端
func (h *RespHandler) Close() error {
	logger.Info("handler shutting down...")
	// 将 closing标志位设置为true，表示Redis服务端正在关闭中
	// 遍历存放客户端连接的activeConn，为其中的每个客户端连接
	// 执行关闭
	h.closing.Set(true)
	h.activeConn.Range(func(key, value any) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		// 返回true 则会继续遍历映射后面的元素，
		// 返回false则不会继续遍历
		return true
	})
	h.db.Close()
	return nil
}
