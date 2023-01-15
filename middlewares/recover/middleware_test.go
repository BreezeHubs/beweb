package recover

import (
	"fmt"
	"github.com/BreezeHubs/beweb"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewMiddlewareBuilder().
		SetPanicResponse(510, []byte("SERVER PANIC")).
		SetLogWithContext(func(ctx *beweb.Context) {
			fmt.Println("LogWithContext: ", string(ctx.ResponseContent))
		}).
		SetLogWithErr(func(err any) {
			fmt.Println("LogWithErr: ", err.(string))
		}).
		SetLogWithContext(func(ctx *beweb.Context) {
			fmt.Println("LogWithContext1: ", string(ctx.ResponseContent))
		}).
		//SetLogWithStack(func(stack string) {
		//	fmt.Println("LogWithStack: ", stack)
		//}).
		Build()

	s := beweb.NewHTTPServer(beweb.WithMiddlewares(builder))
	s.Get("/panic", func(ctx *beweb.Context) {
		value, _ := ctx.QueryValue("dev")
		fmt.Println("value", value)
		if value == "1" {
			panic("panic aaaaaa")
		}
	})

	s.Start(":8080")
}
