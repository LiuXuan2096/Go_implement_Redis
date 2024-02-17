package database

import "go_redis/interface/resp"

// CmdLine 是 [][]byte的别名，表示一个命令行输入
type CmdLine = [][]byte

// Database 接口代表的是Redis业务层的 存储引擎
type Database interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply
	AfterClientClose(c resp.Connection)
	Close()
}

// DataEntity 代表Redis存储的数据的实体
// 将data和key绑定，包括string list has set
type DataEntity struct {
	Data interface{}
}
