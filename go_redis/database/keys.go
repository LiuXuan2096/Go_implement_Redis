package database

import (
	"go_redis/interface/resp"
	"go_redis/lib/wildcard"
	"go_redis/resp/reply"
)

/*
 * 执行Redis中和 key 有关的指令
 */

// execDel 从数据库中移除一个key
func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}

	deleted := db.Removes(keys...)
	return reply.MakeIntReply(int64(deleted))
}

// execExists 检查数据库中是否存在这些key
func execExists(db *DB, args [][]byte) resp.Reply {
	result := int64(0)
	for _, arg := range args {
		key := string(arg)
		_, exists := db.GetEntity(key)
		if exists {
			result++
		}
	}
	return reply.MakeIntReply(result)
}

// execFlushDB 清空当前数据库中的数据
func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	return &reply.OkReply{}
}

// execType 根据key返回数据库中实体的类型
// 包括：string list hash set  zset
// 当前版本只实现了 string类型相关的功能
func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeStatusReply("none")
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
	}
	return &reply.UnKnownErrReply{}
}

// execRename 更改某个key-value中 的key，例如将key1-value 改为 key2-value
// 如果原数据库中已经存在另一个键值对是 key2-value2
// 则会覆盖key2-value2，修改后的数据库中无key1 只有key2-value
func execRename(db *DB, args [][]byte) resp.Reply {
	if len(args) != 2 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'rename' command")
	}
	src := string(args[0])
	dest := string(args[1])

	entity, ok := db.GetEntity(src)
	if !ok {
		return reply.MakeErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)
	return &reply.OkReply{}
}

// execRenameNx : 修改数据库中某个key-value的key
// 只有在新的key 在现有的数据库中不存在时，才会执行修改
func execRenameNx(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])

	_, ok := db.GetEntity(dest)
	if ok {
		return reply.MakeIntReply(0)
	}

	entity, ok := db.GetEntity(src)
	if !ok {
		return reply.MakeErrReply("no such key")
	}
	// 将新key和旧key都从数据库中移除
	// 感觉这里可能会有并发问题
	db.Removes(src, dest)
	db.PutEntity(dest, entity)
	return reply.MakeIntReply(1)
}

// execKeys 根据给定的正则表达式返回 匹配的keys
func execKeys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return reply.MakeMultiBulkReply(result)
}

func init() {
	RegisterCommand("Del", execDel, -2)
	RegisterCommand("Exists", execExists, -2)
	RegisterCommand("Keys", execKeys, 2)
	RegisterCommand("FlushDB", execFlushDB, -1)
	RegisterCommand("Type", execType, 2)
	RegisterCommand("Rename", execRename, 3)
	RegisterCommand("RenameNx", execRenameNx, 3)
}
