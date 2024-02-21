package dict

/*
 * 定义了Redis使用的数据结构的接口，当前版本使用Go官方提供的map实现
 * 后续想要改动，只需要修改接口的实现，不需要变动业务层的代码
 */

// Consumer 用来遍历字典，返回false时会终止遍历
type Consumer func(key string, val interface{}) bool

// Dict 是key-value型数据库的接口
type Dict interface {
	Get(key string) (val interface{}, exists bool)
	Len() int
	Put(key string, val interface{}) (result int)
	PutIfAbsent(key string, val interface{}) (result int)
	PutIfExists(key string, val interface{}) (result int)
	Remove(key string) (result int)
	ForEach(consumer Consumer)
	Keys() []string
	RandomKeys(limit int) []string
	RandomDistinctKeys(limit int) []string
	Clear()
}
