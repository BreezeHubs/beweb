package test

import (
	"fmt"
	"github.com/BreezeHubs/beweb"
	"testing"
)

func TestRouteGroupAndMiddlewares(t *testing.T) {
	m1 := beweb.Middleware(func(next beweb.HandleFunc) beweb.HandleFunc {
		return func(ctx *beweb.Context) {
			fmt.Println("m1")
			next(ctx)
		}
	})
	m2 := beweb.Middleware(func(next beweb.HandleFunc) beweb.HandleFunc {
		return func(ctx *beweb.Context) {
			fmt.Println("m2")
			next(ctx)
		}
	})
	m3 := beweb.Middleware(func(next beweb.HandleFunc) beweb.HandleFunc {
		return func(ctx *beweb.Context) {
			fmt.Println("m3")
			next(ctx)
		}
	})

	s := beweb.NewHTTPServer()

	group := s.Group("/api/v1", m1, m2)
	group.Get("/user", func(ctx *beweb.Context) {
		ctx.Response(200, []byte("user"))
	}, m3)

	s.Get("/name", func(ctx *beweb.Context) {
		ctx.Response(200, []byte("name"))
	})

	s.Start(":8080")
}
