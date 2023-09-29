package tcp

import (
	"bufio"
	"context"
	"errors"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

type EchoHandler struct {
	activeConn sync.Map  //保存，记录有多少连接数
	closing atomic.Boolean
}

// EchoClient is client for EchoHandler, using for test
type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (e *EchoClient) Close() error {
	// 半优雅退出
	e.Waiting.WaitWithTimeout(10*time.Second)
	e.Conn.Close()
	return nil
}

// Handle echos received line to client
func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if handler.closing.Get() {
		// 已经是关闭状态
		// closing handler refuse new connection
		_ = conn.Close()
	}
	// 初始化连接
	client := &EchoClient{
		Conn:    conn,
		Waiting: wait.Wait{},
	}
	handler.activeConn.Store(client, struct{}{})// 把一个客户端存进去
	reader := bufio.NewReader(conn)
	for {
		// may occurs: client EOF, client timeout, server early close
		msg, err := reader.ReadString('\n')
		if err != nil {
			//接收数据出现问题
			if err == io.EOF {
				// 客户端退出
				logger.Info("connection close")
				handler.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Waiting.Add(1)
		time.Sleep(3*time.Second)
		b := []byte(msg)  
		conn.Write(b)  // 回写数据	
	}
}

func (handler *EchoHandler) Close() error {
	if handler.closing.Get() {
		return errors.New("handler is closed!")
	}
	handler.activeConn.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		handler.activeConn.Delete(key)
		return true
	})
	return nil 
}

func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}


