package resp

// Connection 表示一个与redis客户端的连接
type Connection interface {
	Write([]byte) error // 向客户端发送数据
	GetDBIndex() int    // redis内部默认分为16个数据库，返回当前使用的数据库的索引
	SelectDB(int)       // 切换使用的数据库
}

// Reply 是RESP(redis serialization protocol)向客户端发送的消息的接口
type Reply interface {
	ToBytes() []byte
}
