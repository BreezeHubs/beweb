package beweb

import "net/http"

type HandleFunc func(ctx *Context)

// Group 处理路由分组
func (s HTTPServer) Group(groupName string, middlewares ...Middleware) *HTTPServer {
	s.groupNameCache = groupName          //缓存group name
	s.groupMiddlewaresCache = middlewares //缓存middlewares
	return &s
}

// wrapMiddlewares 合并middlewares
func (s *HTTPServer) wrapMiddlewares(middlewares ...Middleware) {
	if len(middlewares) > 0 {
		s.groupMiddlewaresCache = append(s.groupMiddlewaresCache, middlewares...)
	}
}

func (s *HTTPServer) Get(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	s.wrapMiddlewares(middlewares...)
	s.addRoute(s.groupNameCache, http.MethodGet, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Post(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	s.wrapMiddlewares(middlewares...)
	s.addRoute(s.groupNameCache, http.MethodPost, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Delete(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	s.wrapMiddlewares(middlewares...)
	s.addRoute(s.groupNameCache, http.MethodDelete, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Head(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	s.wrapMiddlewares(middlewares...)
	s.addRoute(s.groupNameCache, http.MethodHead, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Options(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	s.wrapMiddlewares(middlewares...)
	s.addRoute(s.groupNameCache, http.MethodOptions, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Put(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	s.wrapMiddlewares(middlewares...)
	s.addRoute(s.groupNameCache, http.MethodPut, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Patch(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	s.wrapMiddlewares(middlewares...)
	s.addRoute(s.groupNameCache, http.MethodPatch, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Connect(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	s.wrapMiddlewares(middlewares...)
	s.addRoute(s.groupNameCache, http.MethodConnect, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Trace(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	s.wrapMiddlewares(middlewares...)
	s.addRoute(s.groupNameCache, http.MethodTrace, path, handleFunc, s.groupMiddlewaresCache...)
}
