package tcp

import (
	"bufio"
	"context"
	"go_redis/lib/logger"
	"go_redis/lib/sync/atomic"
	"go_redis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

/**
 * A echo server to test whether the server is functioning normally
 */

// EchoHandler echos received line to client, using for test
type EchoHandler struct {
	activeConn sync.Map // 使用并发安全的map存储当前有几个客户端连接
	// 标记当前Handler是否已经关闭，如果关闭则不再接受新的客户端连接
	closing atomic.Boolean
}

// EchoClient is client for EchoHandler, using for test.
type EchoClient struct {
	Conn     net.Conn
	Waitting wait.Wait
}

// Handle 将客户端发来的信息，再返回给客户端
func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		// 进入这个if语句说明当前EchoHandler证在准备关闭，此时 拒绝建立新的连接。
		_ = conn.Close()
	}
	client := &EchoClient{
		Conn: conn,
	}
	// 将新的客户端连接 加入到活跃连接集合存储
	h.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn) //使用系统提供的缓存工具，方便我们处理从TCP客户端发来的数据
	for {
		//可能出现的情况：client EOF, client timeout, server early close
		msg, err := reader.ReadString('\n') //设置处理数据的分隔符
		if err != nil {
			if err == io.EOF {
				// 说明客户端断开连接
				logger.Info("connection close")
				h.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Waitting.Add(1) //给WaitGroup的counter加1
		// 将客户端发来的数据再返回给客户端
		b := []byte(msg)
		_, _ = conn.Write(b)

		client.Waitting.Done()
	}
}

// Close 关闭EchoHandler
func (h *EchoHandler) Close() error {
	logger.Info("handler shutting down...")
	h.closing.Set(true) // 将EchoHandler的关闭标志位置为true
	h.activeConn.Range(func(key interface{}, value interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		// 这里返回true是表明，继续遍历后面的元素，如果返回false则不再遍历后面的元素
		// 也即不再关闭后面的客户端连接，目前这个测试功能体现不出这种设计的意义
		// 后面处理更复杂的Redis本身的业务时，会用到这种设计
		return true
	})
	return nil
}

// Close EchoHandler关闭时会调用该方法
func (c *EchoClient) Close() error {
	c.Waitting.WaitWithTimeout(10 * time.Second)
	c.Conn.Close()
	return nil
}

// MakeHandler 新建一个EchoHandler，
func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}
