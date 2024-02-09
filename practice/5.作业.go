package main

import "time"

func main() {
	for i := 0; i < 200000; i++ {
		go func() {
			time.Sleep(5 * time.Second)
		}()
	}

	time.Sleep(1000 * time.Second)
}
