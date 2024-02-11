package main

import (
	"sync/atomic"
)

func add(a *int32) {
	//(*a)++
	atomic.AddInt32(a, 1)
}

//func main() {
//	c := int32(0)
//	for i := 0; i < 1000; i++ {
//		go add(&c)
//	}
//	time.Sleep(1 * time.Second)
//	fmt.Println(c)
//}
