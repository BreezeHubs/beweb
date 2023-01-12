package beweb

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHTTPServer_Middleware(t *testing.T) {
	s := NewHTTPServer()
	s.middlewares = []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第一个before")
				next(ctx)
				fmt.Println("第一个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第二个before")
				next(ctx)
				fmt.Println("第二个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第三个中断")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第四个看不到")
			}
		},
	}

	s.ServeHTTP(nil, &http.Request{})
}
