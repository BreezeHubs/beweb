package beweb

import (
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"unsafe"
)

var json = jsoniter.ConfigFastest

type Context struct {
	Req  *http.Request       //http包的Request
	Resp http.ResponseWriter //http包的ResponseWriter

	PathParams   map[string]string //路由参数缓存
	QueryParams  map[string]string //url get参数缓存
	FormParams   map[string]string //表单参数缓存
	HeaderParams map[string]string //请求头缓存
	Body         bodyBytes         //body缓存

	MatchedRoute string //完整的命中的路由

	ResponseStatus  int
	ResponseContent []byte

	templateEngine TemplateEngine
}

type bodyBytes []byte

func (b bodyBytes) GetBody() string {
	return *(*string)(unsafe.Pointer(&b))
}
