package beweb

import (
	"fmt"
	"github.com/BreezeHubs/beweb/util"
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
	//h.Get("/"+RouteRegExp("order")+"/detail", func(ctx *Context) {
	//	ctx.Resp.Write([]byte("hello world"))
	//})

	h.Get("/orderorder/*", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, " + ctx.Req.URL.Path))
	})
	//h.Get("/test/*", func(ctx *Context) {
	//	ctx.Resp.Write([]byte("hello, " + ctx.Req.URL.Path))
	//})
	h.Get("/test/Reg(^\\d{4}-\\d{8}$)", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, " + ctx.Req.URL.Path))
	})

	//h.Get("/order/*", func(ctx *Context) {
	//	ctx.Resp.Write([]byte("hello, " + ctx.Req.URL.Path))
	//})

	h.Get("/*/abc", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, abc " + ctx.Req.URL.Path))
	})

	h.Get("/*/order", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, order " + ctx.Req.URL.Path))
	})

	h.Get("/param/:id", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, param " + ctx.PathParams["id"]))
	})

	h.Get("/xml/:id", func(ctx *Context) {
		type xml struct {
			Id   int    `xml:"id"`
			Name string `xml:"name"`
		}
		id, _ := ctx.PathParam("id").Int()
		util.ResponseXML(ctx, 200, &xml{
			Id:   id,
			Name: "haha",
		})
	})

	h.Start(":8080")
}
