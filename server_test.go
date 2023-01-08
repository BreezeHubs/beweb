package beweb

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	h := NewHTTPServer()

	h.Get("/user", func(ctx *Context) {
		fmt.Println("1")
		fmt.Println("2")
	})

	h.Get("/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello world"))
	})

	h.Get("/order/*", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, " + ctx.Req.URL.Path))
	})

	h.Get("/*/abc", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, abc " + ctx.Req.URL.Path))
	})

	h.Get("/*/order", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, order " + ctx.Req.URL.Path))
	})

	h.Get("/param/:id", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, param " + ctx.PathParams["id"]))
	})

	h.Start(":8080")
}
