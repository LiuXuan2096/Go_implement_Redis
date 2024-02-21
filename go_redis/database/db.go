package database

import (
	"go_redis/datastructure/dict"
	"go_redis/interface/database"
	"go_redis/interface/resp"
	"go_redis/resp/reply"
	"strings"
)

// DB 1.存储数据 2.执行用户指令
type DB struct {
	index int
	// key -> DataEntity 键值对
	data dict.Dict
}

// ExecFunc 是用户命令的executor的接口
type ExecFunc func(db *DB, args [][]byte) resp.Reply

// CmdLine 是[][]byte的别名，表示用户通过客户端传来的一条指令
type CmdLine = [][]byte

// makeDB 创建一个 DB 实例
func makeDB() *DB {
	db := &DB{
		data: dict.MakeSycnDict(),
	}
	return db
}

// validateArity 进行参数个数校验，判断传来的指令的参数个数是否合法
func validateArity(arity int, cmdArgs CmdLine) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		return arity == argNum
	}
	// 有些指令的参数数目不确定 比如批量删除数据库中一些key的指令 和 ping指令
	// Del key1 key2 key3 ......
	// 这种类型指令的arity设置为它们最少应该具有的参数数量的负值：-2 或 -1
	return argNum >= -arity
}

// Exec 在一个数据库上执行指令
func (db *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {

	cmdName := strings.ToLower(string(cmdLine[0])) // 全部转成小写
	// 从cmdTable中根据用户传来的指令的第一个单词，取出对应的指令的执行函数
	cmd, ok := cmdTable[cmdName]
	if !ok {
		// 说明cmdTable中没有该指令的执行器
		return reply.MakeErrReply("ERR unknow command '" + cmdName + "'")
	}
	if !validateArity(cmd.arity, cmdLine) {
		return reply.MakeArgNumErrReply(cmdName)
	}
	return cmd.executor(db, cmdLine[1:])
}

/* -------- 数据访问 ------- */

// GetEntity 返回数据库中key对应的 DataEntity
func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	entity, _ := raw.(*database.DataEntity)
	return entity, true
}

// PutEntity 向数据库中插入一个key-value
// 如果数据库中已经存在key 会覆盖；如果不存在 会新增一个key-value
func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

// PutIfExists 更新一个已经存在的key-value
// 如果该key不存在 返回0
func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

// PutIfAbsent 只有在该key不存在时才会向数据库新增一个key-value
// 如果该key已经存在，则不执行任何操作 返回0
func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

// Remove 从数据库中移除指定的key
func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

// Removes 将给定的key全部从数据库中移除
func (db *DB) Removes(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.Remove(key)
			deleted++
		}
	}
	return deleted
}

// Flush 清空数据库
func (db *DB) Flush() {
	db.data.Clear()
}
