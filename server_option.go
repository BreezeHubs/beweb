package beweb

type HTTPServerOpt func(*HTTPServer)

// WithGracefullyExit 优雅退出设置
func WithGracefullyExit(b bool, fn func()) HTTPServerOpt {
	return func(s *HTTPServer) {
		s.isGracefullyExit = b
		s.isGracefullyExitFunc = fn
	}
}
