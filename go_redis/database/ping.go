package database

import (
	"go_redis/interface/resp"
	"go_redis/resp/reply"
)

// 处理和 ping 指令有关的逻辑

// Ping 用户ping Redis 服务端时，向客户端返回 Pong
func Ping(db *DB, args [][]byte) resp.Reply {
	if len(args) == 0 {
		return &reply.PongReply{}
	} else if len(args) == 1 {
		return reply.MakeStatusReply(string(args[0]))
	} else {
		return reply.MakeErrReply("ERR wrong number of arguments for 'ping' command")
	}
}

func init() {
	RegisterCommand("ping", Ping, -1)
}
