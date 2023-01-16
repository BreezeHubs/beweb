package beweb

import (
	"context"
	"fmt"
	"github.com/BreezeHubs/bekit/sys"
	"net/http"
	"time"
)

func NewHTTPServer(opts ...HTTPServerOpt) IHTTPServer {
	s := &HTTPServer{
		router: newRouter(), //路由树

		//默认配置
		gracefullyExitTimeout: 30 * time.Second,
		shutdownTimeout:       10 * time.Second,
	}

	//执行配置函数
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start 运行web服务
func (s *HTTPServer) Start(addr string) error {
	//创建原生 http server
	server := &http.Server{Addr: addr, Handler: s}
	go func() {
		//run http server
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	//是否开启优雅退出
	if s.isGracefullyExit {
		//创建退出信号监听
		signal := sys.NewListenExitSignal()

		//监听退出信号
		for !signal.IsExit() {
			//fmt.Println("running...")
			time.Sleep(10 * time.Millisecond)
		}

		//创建超时控制
		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		err := server.Shutdown(ctx)

		//channel监听超时
		done := make(chan struct{}, 1)
		go func() {
			//执行自定义回收
			s.isGracefullyExitFunc()
			done <- struct{}{}
		}()

		select {
		case <-done:
		case <-time.After(s.gracefullyExitTimeout):
			fmt.Println("beweb gracefully timeout")
		}

		fmt.Println("beweb exit")
		return err
	} else {
		select {}
	}
	return nil
}
