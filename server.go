package beweb

import (
	"fmt"
	"github.com/BreezeHubs/bekit/sys"
	"net/http"
	"time"
)

type HandleFunc func(ctx *Context)

func NewHTTPServer() IHTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

// Start 运行web服务
func (s *HTTPServer) Start(addr string) error {
	//创建退出信号监听
	signal := sys.NewListenExitSignal()

	go http.ListenAndServe(addr, s) //http server

	//监听退出信号
	for !signal.IsExit() {
		fmt.Println("running...")
		time.Sleep(1 * time.Second)
	}

	//进行一些回收动作...
	fmt.Println("test：进行一些回收动作...")

	fmt.Println("exit")
	return nil
}
