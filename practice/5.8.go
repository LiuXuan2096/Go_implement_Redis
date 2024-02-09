package main

import (
	"fmt"
	"time"
)

func do(i int, ch chan struct{}) {
	fmt.Println("葛诗颖", i)
	time.Sleep(time.Second)
	<-ch
}

//func main() {
//	c := make(chan struct{}, 3000)
//	for i := 0; i < math.MaxInt32; i++ {
//		c <- struct{}{}
//		go do(i, c)
//	}
//	time.Sleep(time.Hour)
//}
