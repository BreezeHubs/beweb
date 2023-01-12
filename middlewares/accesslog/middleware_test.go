package accesslog

import (
	"fmt"
	"github.com/BreezeHubs/beweb"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	mdl := NewMiddlewareBuilder().
		LogInputFunc(func(ctx *beweb.Context) (string, error) {
			l := struct {
				ServerName string `json:"server_name"`
				Host       string `json:"host"`
				Route      string `json:"route"` //完整的命中的路由
				HTTPMethod string `json:"http_method"`
				Path       string `json:"path"`
			}{
				ServerName: "test-server",
				Host:       ctx.Req.Host,
				Route:      ctx.MatchedRoute, //完整的命中的路由
				HTTPMethod: ctx.Req.Method,
				Path:       ctx.Req.URL.Path,
			}
			data, err := json.Marshal(l)
			return string(data), err
		}).
		LogOutFunc(func(logString string, err error) {
			fmt.Println(logString)
		}).Build()

	s := beweb.NewHTTPServer(
		beweb.WithMiddlewares(mdl),
	)

	s.Post("/a/*/c", func(ctx *beweb.Context) {
		fmt.Println("hello")
	})
	req, _ := http.NewRequest(http.MethodPost, "/a/b/c", nil)
	req.Host = "127.0.0.1"
	s.ServeHTTP(nil, req)
}
