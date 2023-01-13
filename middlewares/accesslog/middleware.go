package accesslog

import (
	"github.com/BreezeHubs/beweb"
	jsoniter "github.com/json-iterator/go"
	"log"
	"time"
	"unsafe"
)

var json = jsoniter.ConfigFastest

type MiddlewareBuilder struct {
	logInputFunc func(ctx *beweb.Context) (string, error) //log的输出方式
	logOutFunc   func(logString string, err error)        //log的输出格式
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		//初始化log的输出方式
		logInputFunc: func(ctx *beweb.Context) (string, error) {
			// 默认log数据
			type accessLog struct {
				Time       string `json:"time"`
				Host       string `json:"host"`
				Route      string `json:"route"` //命中的路由
				HTTPMethod string `json:"http_method"`
				Path       string `json:"path"`
				//Header     map[string]string `json:"header"`
				//Body       string             `json:"body"`
				//Response   beweb.ResponseData `json:"response"`
			}

			l := accessLog{
				Time:       time.Now().Format("2006-01-02 15:04:05.999999999"),
				Host:       ctx.Req.Host,
				Route:      ctx.MatchedRoute, //完整的命中的路由
				HTTPMethod: ctx.Req.Method,
				Path:       ctx.Req.URL.Path,
				//Header:     ctx.HeaderParams,
				//Body:       ctx.Body.GetBody(),
				//Response:   ctx.ResponseData,
			}

			data, err := json.Marshal(l)
			return *(*string)(unsafe.Pointer(&data)), err
		},
		//初始化log的输出格式
		logOutFunc: func(logString string, err error) {
			log.Println(logString)
		},
	}
}

// LogOutFunc 定义log的输出方式
func (m *MiddlewareBuilder) LogOutFunc(fn func(logString string, err error)) *MiddlewareBuilder {
	m.logOutFunc = fn
	return m
}

// LogInputFunc 定义log的输出格式
func (m *MiddlewareBuilder) LogInputFunc(fn func(ctx *beweb.Context) (string, error)) *MiddlewareBuilder {
	m.logInputFunc = fn
	return m
}

// Build 构建
func (m *MiddlewareBuilder) Build() beweb.Middleware {
	return func(next beweb.HandleFunc) beweb.HandleFunc {
		return func(ctx *beweb.Context) {
			//记录请求
			defer func() {
				data, err := m.logInputFunc(ctx) //log 输入
				m.logOutFunc(data, err)          //log 输出
			}()
			next(ctx)
		}
	}
}
