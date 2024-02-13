package main

//import (
//	"fmt"
//	"time"
//)

//func main() {
//
//	go func() {
//
//		panic("")
//		fmt.Println("end g")
//	}()
//
//	time.Sleep(time.Second)
//	fmt.Println("end main")
//}

//func main() {
//	defer fmt.Println("main g")
//	go func() {
//		defer fmt.Println("defer g")
//		panic("")
//		fmt.Println("end g")
//	}()
//
//	time.Sleep(time.Second)
//	fmt.Println("end main")
//}

//func main() {
//	defer fmt.Println("main g")
//	go func() {
//		defer func() {
//			recover()
//		}()
//		panic("")
//		fmt.Println("end g")
//	}()
//
//	time.Sleep(time.Second)
//	fmt.Println("end main")
//}
