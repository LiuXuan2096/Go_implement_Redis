package main

import "sync"

/*
使用 $env:GODEBUG="gctrace=1" 查看gc情况
*/

func main() {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup) {
			var count int
			for i := 0; i < 1e10; i++ {
				count++
			}
			wg.Done()
		}(&wg)
	}
}
