package tcp

import (
	"context"
	"fmt"
	"go-redis/interface/tcp"
	"go-redis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Config stores tcp server properties
type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeCh := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeCh <- struct{}{}
		}
	}()
	
	lis, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("bind: %s, start listening...", cfg.Address))
	// 调用ListenAndServe
	return ListenAndServe(lis, handler, closeCh)

}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeCh <-chan struct{}) error {
	go func() {  // 监听关闭关闭 通道
		<-closeCh
		logger.Info("shutting down")
		_ = listener.Close()
		_ = handler.Close()
	}()
	defer func() {
		// 保证连接关闭
		listener.Close()
		handler.Close()
	}()
	ctx := context.Background()
	var wg sync.WaitGroup
	var err error
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accepted link")
		wg.Add(1)  // 服务 +1
		go func() {
			defer func() {
				wg.Done()  // 确保服务 -1
			}()
			// 处理handler
			handler.Handle(ctx, conn)
		}()
	}
	wg.Wait()
	return err
}
