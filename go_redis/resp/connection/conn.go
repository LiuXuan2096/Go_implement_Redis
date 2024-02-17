package connection

import (
	"go_redis/lib/sync/wait"
	"net"
	"sync"
	"time"
)

// Connection 代表的是 一个同Redis客户端的连接
type Connection struct {
	conn net.Conn
	// 等到向客户端发送消息的协程都结束工作后，再关闭同Redis客户端的连接
	waitingReply wait.Wait
	// 在handler向客户端发送消息时上锁，
	// 确保同时只能有一个Connection向客户端发送数据
	mu sync.Mutex
	// 表示选择的数据库引擎的索引
	selectedDB int
}

// NewConn 建立一个新的同Redis客户端的连接
func NewConn(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

// RemoteAddr 返回连接的客户端的网络地址
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// Close 断开同客户端的连接
func (c *Connection) Close() error {
	c.waitingReply.WaitWithTimeout(10 * time.Second)
	_ = c.conn.Close()
	return nil
}

// Write 通过TCP 连接向客户端发送数据
func (c *Connection) Write(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	c.mu.Lock()
	c.waitingReply.Add(1)
	defer func() {
		c.waitingReply.Done()
		c.mu.Unlock()
	}()

	_, err := c.conn.Write(b)
	return err
}

// GetDBIndex 返回当前使用的数据库的索引
func (c *Connection) GetDBIndex() int {
	return c.selectedDB
}

// SelectDB 选择一个数据库
func (c *Connection) SelectDB(dbNum int) {
	c.selectedDB = dbNum
}
