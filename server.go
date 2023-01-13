package beweb

import (
	"context"
	"fmt"
	"github.com/BreezeHubs/bekit/sys"
	"net/http"
	"time"
)

type HandleFunc func(ctx *Context)

func NewHTTPServer(opts ...HTTPServerOpt) *HTTPServer {
	s := &HTTPServer{
		router:                newRouter(),
		gracefullyExitTimeout: 30 * time.Second,
		shoutdownTimeout:      10 * time.Second,
	}

	//执行配置
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start 运行web服务
func (s *HTTPServer) Start(addr string) error {
	//创建退出信号监听
	signal := sys.NewListenExitSignal()

	//创建http server
	server := &http.Server{Addr: addr, Handler: s}
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.ListenAndServe() //run http server
	}()

	//监听退出信号
	for !signal.IsExit() {
		//fmt.Println("running...")
		time.Sleep(10 * time.Millisecond)
	}

	//进行一些回收动作...
	if s.isGracefullyExit {
		ctx, cancel := context.WithTimeout(context.Background(), s.shoutdownTimeout)
		defer cancel()
		err := server.Shutdown(ctx)

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

		fmt.Println("beweb gracefully exit")
		return err
	}

	err := <-errChan
	return err
}
