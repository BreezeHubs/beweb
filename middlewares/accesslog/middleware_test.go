package accesslog

import (
	"fmt"
	"github.com/BreezeHubs/beweb"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{}
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()

	s := beweb.NewHTTPServer(
		beweb.WithMiddlewares(mdl),
	)

	s.Post("/a/b/c", func(ctx *beweb.Context) {
		fmt.Println("hello")
	})
	req, _ := http.NewRequest(http.MethodPost, "/a/b/c", nil)
	s.ServeHTTP(nil, req)
}
