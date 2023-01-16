package beweb

import "net/http"

type IHTTPServer interface {
	http.Handler

	Group(groupName string, middlewares ...Middleware) *HTTPServer
	/*
		请求方法
	*/
	Get(path string, handleFunc HandleFunc, middlewares ...Middleware)
	Post(path string, handleFunc HandleFunc, middlewares ...Middleware)
	Delete(path string, handleFunc HandleFunc, middlewares ...Middleware)
	Head(path string, handleFunc HandleFunc, middlewares ...Middleware)
	Options(path string, handleFunc HandleFunc, middlewares ...Middleware)
	Put(path string, handleFunc HandleFunc, middlewares ...Middleware)
	Patch(path string, handleFunc HandleFunc, middlewares ...Middleware)
	Connect(path string, handleFunc HandleFunc, middlewares ...Middleware)
	Trace(path string, handleFunc HandleFunc, middlewares ...Middleware)

	// Start 运行web服务
	Start(addr string) error
}
