package beweb

import "net/http"

func (s HTTPServer) Group(groupName string, middlewares ...Middleware) *HTTPServer {
	s.groupNameCache = groupName
	s.groupMiddlewaresCache = middlewares
	return &s
}

func (s *HTTPServer) Get(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	if len(middlewares) > 0 {
		s.groupMiddlewaresCache = append(s.groupMiddlewaresCache, middlewares...)
	}
	s.addRoute(s.groupNameCache, http.MethodGet, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Post(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	if len(middlewares) > 0 {
		s.groupMiddlewaresCache = append(s.groupMiddlewaresCache, middlewares...)
	}
	s.addRoute(s.groupNameCache, http.MethodPost, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Delete(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	if len(middlewares) > 0 {
		s.groupMiddlewaresCache = append(s.groupMiddlewaresCache, middlewares...)
	}
	s.addRoute(s.groupNameCache, http.MethodDelete, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Head(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	if len(middlewares) > 0 {
		s.groupMiddlewaresCache = append(s.groupMiddlewaresCache, middlewares...)
	}
	s.addRoute(s.groupNameCache, http.MethodHead, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Options(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	if len(middlewares) > 0 {
		s.groupMiddlewaresCache = append(s.groupMiddlewaresCache, middlewares...)
	}
	s.addRoute(s.groupNameCache, http.MethodOptions, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Put(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	if len(middlewares) > 0 {
		s.groupMiddlewaresCache = append(s.groupMiddlewaresCache, middlewares...)
	}
	s.addRoute(s.groupNameCache, http.MethodPut, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Patch(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	if len(middlewares) > 0 {
		s.groupMiddlewaresCache = append(s.groupMiddlewaresCache, middlewares...)
	}
	s.addRoute(s.groupNameCache, http.MethodPatch, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Connect(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	if len(middlewares) > 0 {
		s.groupMiddlewaresCache = append(s.groupMiddlewaresCache, middlewares...)
	}
	s.addRoute(s.groupNameCache, http.MethodConnect, path, handleFunc, s.groupMiddlewaresCache...)
}

func (s *HTTPServer) Trace(path string, handleFunc HandleFunc, middlewares ...Middleware) {
	if len(middlewares) > 0 {
		s.groupMiddlewaresCache = append(s.groupMiddlewaresCache, middlewares...)
	}
	s.addRoute(s.groupNameCache, http.MethodTrace, path, handleFunc, s.groupMiddlewaresCache...)
}
