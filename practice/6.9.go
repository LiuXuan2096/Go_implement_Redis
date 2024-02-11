package main

//import (
//	"fmt"
//	"sync"
//	"time"
//)
//
//func GeShiYin(w *sync.WaitGroup) {
//	fmt.Println("葛诗颖在踢足球")
//	w.Done()
//}
//
//func main() {
//	wg := sync.WaitGroup{}
//	wg.Add(3)
//	go GeShiYin(&wg)
//	go GeShiYin(&wg)
//	go func(wg *sync.WaitGroup) {
//		fmt.Println("葛诗颖")
//		time.Sleep(time.Second)
//		wg.Done()
//	}(&wg)
//	wg.Wait()
//}
