package beweb

import "net/http"

type IHTTPServer interface {
	http.Handler

	Get(path string, handleFunc HandleFunc)
	Post(path string, handleFunc HandleFunc)
	Delete(path string, handleFunc HandleFunc)
	Head(path string, handleFunc HandleFunc)
	Options(path string, handleFunc HandleFunc)
	Put(path string, handleFunc HandleFunc)
	Patch(path string, handleFunc HandleFunc)

	// Start 运行web服务
	Start(addr string) error
}
