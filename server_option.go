package beweb

type HTTPServerOpt func(*HTTPServer)

func WithGracefullyExit(b bool, fn func()) HTTPServerOpt {
	return func(s *HTTPServer) {
		s.isGracefullyExit = b
		s.isGracefullyExitFunc = fn
	}
}
