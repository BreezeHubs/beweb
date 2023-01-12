package accesslog

import (
	"github.com/BreezeHubs/beweb"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type MiddlewareBuilder struct {
	logFunc func(log string)
}

func (m *MiddlewareBuilder) LogFunc(fn func(log string)) *MiddlewareBuilder {
	m.logFunc = fn
	return m
}

func (m MiddlewareBuilder) Build() beweb.Middleware {
	return func(next beweb.HandleFunc) beweb.HandleFunc {
		return func(ctx *beweb.Context) {
			//记录请求
			defer func() {
				l := accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchedRoute,
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}
				data, _ := json.Marshal(l)
				m.logFunc(string(data))
			}()
			next(ctx)
		}
	}
}

type accessLog struct {
	Host       string `json:"host,omitempty"`
	Route      string `json:"route,omitempty"` //命中的路由
	HTTPMethod string `json:"http_method,omitempty"`
	Path       string `json:"path,omitempty"`
}
