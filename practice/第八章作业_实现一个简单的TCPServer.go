package main

//import (
//	"bufio"
//	"fmt"
//	"io"
//	"log"
//	"net"
//)
//
///*
//实现一个简单的TCP Server
//*/
//
//func ListenAndServe(address string) {
//	// 绑定监听地址
//	listen, err := net.Listen("tcp", address)
//	if err != nil {
//		log.Fatal(fmt.Sprintf("listen err: %v", err))
//	}
//	defer listen.Close()
//	log.Println(fmt.Sprintf("bind: %s, start listening...", address))
//
//	for true {
//		// Accept会一直阻塞直到有新的连接建立或者Listen中断才返回
//		conn, err := listen.Accept()
//		if err != nil {
//			// 通常是由于listen被关闭导致无法继续监听导致的错误
//			log.Fatal(fmt.Sprintf("accept err: %v", err))
//		}
//		// 开启新的Goroutine处理该连接
//		go Handle(conn)
//	}
//}
//
//func Handle(conn net.Conn) {
//	// 使用bufio标准库提供的缓冲区功能
//	reader := bufio.NewReader(conn)
//	for {
//		// ReadString 会一直阻塞直到遇到分隔符'\n'
//		// 遇到分隔符会返回分隔符或连接建立后收到的所有数据，包括分隔符本身
//		// 若在遇到分隔符之前遇到异常，ReadString会返回已收到的数据和错误信息
//		msg, err := reader.ReadString('\n')
//		if err != nil {
//			// 通常遇到的错误是连接中断或被关闭，用io.EOF表示
//			if err == io.EOF {
//				log.Println("connection close")
//			} else {
//				log.Println(err)
//			}
//			return
//		}
//		b := []byte(msg)
//		// 将收到的信息发送给客户端
//		conn.Write(b)
//	}
//}
//
//func main() {
//	ListenAndServe(":8080")
//}
