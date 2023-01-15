package accesslog

import (
	"fmt"
	"github.com/BreezeHubs/beweb"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	mdl := NewMiddlewareBuilder().
		//LogInputFunc(func(ctx *beweb.Context) (string, error) {
		//	l := struct {
		//		ServerName string `json:"server_name"`
		//		Host       string `json:"host"`
		//		Route      string `json:"route"`
		//		HTTPMethod string `json:"http_method"`
		//		Path       string `json:"path"`
		//		Time       string `json:"time"`
		//	}{
		//		ServerName: "test-server",
		//		Host:       ctx.Req.Host,
		//		Route:      ctx.MatchedRoute,
		//		HTTPMethod: ctx.Req.Method,
		//		Path:       ctx.Req.URL.Path,
		//		Time:       time.Now().Format("2006-01-02 15:04:05.999999999"),
		//	}
		//	data, err := json.Marshal(l)
		//	return string(data), err
		//}).
		LogOutFunc(func(logString string, err error) {
			fmt.Println(logString)
		}).Build()

	s := beweb.NewHTTPServer(
		beweb.WithMiddlewares(mdl),
	)

	s.Post("/a/*/c", func(ctx *beweb.Context) {
		fmt.Println("hello")
		ctx.Response(200, []byte(`{"Id":1,"Name":"breeze"}`))
		//util.ResponseJSON(ctx, 200, "", "", struct {
		//	Id   int    `json:"id"`
		//	Name string `json:"name"`
		//}{
		//	Id:   1,
		//	Name: "breeze",
		//})
	})
	req, _ := http.NewRequest(http.MethodPost, "/a/b/c", nil)
	req.Host = "127.0.0.1"
	req.Header.Set("Content-ID", "123456")
	req.Header.Set("Content-Language", "ZH")
	req.Body = io.NopCloser(strings.NewReader("{\"data\":\"hahaha\"}"))
	s.ServeHTTP(nil, req)
}
