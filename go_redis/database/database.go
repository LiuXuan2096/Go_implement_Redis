package database

import (
	"fmt"
	"go_redis/config"
	"go_redis/interface/resp"
	"go_redis/lib/logger"
	"go_redis/resp/reply"
	"runtime/debug"
	"strconv"
	"strings"
)

// Database 是多个 DB 的集合,
// 参考Redis官方的设计，每个 Database 默认有16个 DB
type Database struct {
	dbSet []*DB
}

// NewDatabase 创建一个Redis Database
func NewDatabase() *Database {
	mdb := &Database{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	mdb.dbSet = make([]*DB, config.Properties.Databases)
	for i := range mdb.dbSet {
		singleDB := makeDB()
		singleDB.index = i
		mdb.dbSet[i] = singleDB
	}
	return mdb
}

// Exec 执行客户端发来的Redis指令
// 参数 `cmdLine` 包括了命令和它的参数，例如："set key value"
func (mdb *Database) Exec(c resp.Connection, cmdLine CmdLine) (result resp.Reply) {
	defer func() {
		// 执行用户发来的指令时可能出错，在defer中用recover处理panic
		// 以免发生的错误层层上发 最后带崩整个程序
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
		}
	}()

	cmdName := strings.ToLower(string(cmdLine[0])) // 获取指令类型
	if cmdName == "select" {
		// 切换数据库的指令
		if len(cmdLine) != 2 {
			return reply.MakeArgNumErrReply("select")
		}
		return execSelect(c, mdb, cmdLine[1:])
	}
	// 执行其他的Redis指令时 由 Exec 执行
	dbIndex := c.GetDBIndex()
	selectDB := mdb.dbSet[dbIndex]
	return selectDB.Exec(c, cmdLine)
}

// Close 关闭数据库时，执行的逻辑
func (mdb *Database) Close() {
	// Todo: 当前版本未实现该方法
}

// AfterClientClose 关闭一个同数据库连接的客户端连接后 要执行的逻辑
func (mdb *Database) AfterClientClose(c resp.Connection) {
	// Todo: 当前版本未实现该方法
}

// execSelect Redis中的 切换数据库的 select语句的执行函数
func execSelect(c resp.Connection, mdb *Database, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("ERR invalid DB index")
	}
	if dbIndex >= len(mdb.dbSet) {
		return reply.MakeErrReply("ERR DB index is out of range")
	}
	c.SelectDB(dbIndex)
	return reply.MakeOkReply()
}
