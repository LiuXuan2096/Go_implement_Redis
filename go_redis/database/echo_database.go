package database

import (
	"go_redis/interface/resp"
	"go_redis/lib/logger"
	"go_redis/resp/reply"
)

type EchoDatabase struct {
}

func (e EchoDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	return reply.MakeMultiBulkReply(args)
}

func (e EchoDatabase) AfterClientClose(c resp.Connection) {
	logger.Info("EchoDatabase AfterClientClose")
}

func (e EchoDatabase) Close() {
	logger.Info("EchoDatabase Close")
}

func NewEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}
