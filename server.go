package beweb

import (
	"fmt"
	"github.com/BreezeHubs/bekit/sys"
	"net/http"
	"time"
)

type HandleFunc func(ctx *Context)

func NewHTTPServer(opts ...HTTPServerOpt) *HTTPServer {
	s := &HTTPServer{
		router: newRouter(),
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

	go http.ListenAndServe(addr, s) //http server

	//监听退出信号
	for !signal.IsExit() {
		//fmt.Println("running...")
		time.Sleep(10 * time.Millisecond)
	}

	//进行一些回收动作...
	if s.isGracefullyExit {
		s.isGracefullyExitFunc()
	}

	fmt.Println("exit")
	return nil
}
