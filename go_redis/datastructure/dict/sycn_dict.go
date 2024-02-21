package dict

import "sync"

type SyncDict struct {
	m sync.Map
}

// MakeSycnDict 返回一个新的SyncDict变量
func MakeSycnDict() *SyncDict {
	return &SyncDict{}
}

// Get 返回key对应的value以及 该key是否存在
func (dict *SyncDict) Get(key string) (val interface{}, exists bool) {
	val, ok := dict.m.Load(key)
	return val, ok
}

// Len 返回dict中的元素数量
func (dict *SyncDict) Len() int {
	length := 0
	dict.m.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	return length
}

// Put 向dict中添加key-value键值对，并返回新插入的key-value数量
func (dict *SyncDict) Put(key string, val interface{}) (result int) {
	_, existed := dict.m.Load(key)
	dict.m.Store(key, val)
	if existed {
		return 0 // 说明原dict中就存在key这个键，所以新插入的key-value键值对数量为0
	}
	return 1
}

// PutIfAbsent 当key不存在时，才向dict中添加key-value键值对，并返回更新的key-value的数量
func (dict *SyncDict) PutIfAbsent(key string, val interface{}) (result int) {
	_, existed := dict.m.Load(key)
	if existed {
		return 0
	}
	dict.m.Store(key, val)
	return 1
}

// PutIfExists 当dict中本来就存在key时，才将key-value插入，并返回插入的key-value的数量
func (dict *SyncDict) PutIfExists(key string, val interface{}) (result int) {
	_, existed := dict.m.Load(key)
	if existed {
		dict.m.Store(key, val)
		return 1
	}
	return 0
}

// Remove 移除key对应的key-value，并返回删除的key-value的数量
func (dict *SyncDict) Remove(key string) (result int) {
	_, existed := dict.m.Load(key)
	dict.m.Delete(key)
	if existed {
		return 1
	}
	return 0
}

// ForEach 遍历整个dict，对dict中的每个key-value执行consumer方法
func (dict *SyncDict) ForEach(consumer Consumer) {
	dict.m.Range(func(key, value any) bool {
		consumer(key.(string), value)
		return true
	})
}

// Keys 返回dict中所有的key组成的Slice
func (dict *SyncDict) Keys() []string {
	result := make([]string, dict.Len())
	i := 0
	dict.m.Range(func(key, value any) bool {
		result[i] = key.(string)
		i++
		return true
	})
	return result
}

// RandomKeys 随机返回给定数量的key组成的Slice，可能包含重复的key
func (dict *SyncDict) RandomKeys(limit int) []string {
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		dict.m.Range(func(key, value any) bool {
			result[i] = key.(string)
			return false
		})
	}
	return result
}

// RandomDistinctKeys 随机返回给定数量的key组成的Slice，不会包含重复的key
func (dict *SyncDict) RandomDistinctKeys(limit int) []string {
	result := make([]string, limit)
	i := 0
	dict.m.Range(func(key, value any) bool {
		if i == limit {
			return false
		}
		result[i] = key.(string)
		i++
		return true
	})
	return result
}

// Clear 将dict中的数据清空
func (dict *SyncDict) Clear() {
	*dict = *MakeSycnDict()
}
