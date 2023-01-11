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

func (s *HTTPServer) Get(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodGet, path, handleFunc)
}

func (s *HTTPServer) Post(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodPost, path, handleFunc)
}

func (s *HTTPServer) Delete(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodDelete, path, handleFunc)
}

func (s *HTTPServer) Head(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodHead, path, handleFunc)
}

func (s *HTTPServer) Options(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodOptions, path, handleFunc)
}

func (s *HTTPServer) Put(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodPut, path, handleFunc)
}

func (s *HTTPServer) Patch(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodPatch, path, handleFunc)
}

// Start 运行web服务
func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}
