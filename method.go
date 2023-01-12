package beweb

import "net/http"

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

func (s *HTTPServer) Connect(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodConnect, path, handleFunc)
}

func (s *HTTPServer) Trace(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodTrace, path, handleFunc)
}
