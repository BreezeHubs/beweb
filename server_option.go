package beweb

type HTTPServerOpt func(*HTTPServer)

// WithGracefullyExit 优雅退出设置
func WithGracefullyExit(b bool, fn func()) HTTPServerOpt {
	return func(s *HTTPServer) {
		s.isGracefullyExit = b
		s.isGracefullyExitFunc = fn
	}
}

// WithMiddlewares 添加中间件
func WithMiddlewares(fn ...Middleware) HTTPServerOpt {
	return func(s *HTTPServer) {
		s.middlewares = append(s.middlewares, fn...)
	}
}
