package main

/*
使用Channel开发一个带有非阻塞功能的锁
*/

// 声明MyMutex结构体
type MyMutex chan struct{}

// 构造方法：使用一个缓冲大小为1的Channel，载体为空结构体
func NewMyMutex() MyMutex {
	ch := make(chan struct{}, 1)
	return ch
}

// 加锁时，向Channel塞入一个数据
// 如果已经被加锁，后面的协程无法塞入数据，阻塞。
func (m *MyMutex) Lock() {
	(*m) <- struct{}{}
}
