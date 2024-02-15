package tcp

import (
	"context"
	"fmt"
	"go_redis/interface/tcp"
	"go_redis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

/**
 * A tcp server
 */

// Config stores tcp server properties
type Config struct {
	Address string
}

// ListenAndServeWithSignal binds port and handle requests, blocking until receive stop signal
func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	// 起到发送关闭信号的作用，在程序被关闭时即收到系统发来的关闭信号后，向ListenAndServe方法发送
	// 关闭信号，空结构体即起到发送信号的作用。在ListenAndServe处理程序关闭时具体的善后逻辑
	closeChan := make(chan struct{})
	sigCh := make(chan os.Signal) // 接收系统发来的信号
	// 注册系统要接收的系统信号
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	// 启动一个协程监听系统发来的信号，一旦收到系统发来的关闭程序的信号，就像closeChan发送空结构体作为程序
	// 退出的信号
	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	// 通过Go的net包创建一个监听cfg.Address端口的TCP Socket
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("bind: %s, start listening...", cfg.Address))
	ListenAndServe(listener, handler, closeChan)
	return nil
}

// ListenAndServe binds port and handle requests, blocking until close
func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	// 开启一个协程，监听关闭信号
	go func() {
		<-closeChan
		logger.Info("shutting down...")
		_ = listener.Close()
		_ = handler.Close()
	}()

	// 监听端口

	defer func() {
		// 发生意料外的错误时，关闭TCP连接
		_ = listener.Close()
		_ = handler.Close()
	}()
	ctx := context.Background()
	var waitDone sync.WaitGroup
	for {
		// 通过一个死循环，一直监听连接请求，当收到新的连接请求并成功建立连接时
		// 新建一个协程处理连接后的具体业务逻辑
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		// 新建一个协程，处理连接成功后的具体业务逻辑
		logger.Info("accept link")
		waitDone.Add(1)
		go func() {
			defer waitDone.Done()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
