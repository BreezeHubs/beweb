package beweb

import "time"

type HTTPServerOpt func(*HTTPServer)

// WithGracefullyExit 优雅退出设置
func WithGracefullyExit(b bool, fn func(), timeout ...time.Duration) HTTPServerOpt {
	return func(s *HTTPServer) {
		s.isGracefullyExit = b
		s.isGracefullyExitFunc = fn

		if len(timeout) > 0 {
			s.gracefullyExitTimeout = timeout[0]
		}
	}
}

// WithMiddlewares 添加中间件
func WithMiddlewares(fn ...Middleware) HTTPServerOpt {
	return func(s *HTTPServer) {
		s.middlewares = append(s.middlewares, fn...)
	}
}

// WithShutdownTimeout 设置http退出超时
func WithShutdownTimeout(timeout time.Duration) HTTPServerOpt {
	return func(s *HTTPServer) {
		s.shutdownTimeout = timeout
	}
}
