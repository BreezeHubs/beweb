package beweb

import (
	"net/http"
)

type HandleFunc func(ctx *Context)

func NewHTTPServer() IHTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

// Start 运行web服务
func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}
