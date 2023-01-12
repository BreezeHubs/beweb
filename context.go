package beweb

import (
	"net/http"
)

type Context struct {
	Req  *http.Request       //http包的Request
	Resp http.ResponseWriter //http包的ResponseWriter

	PathParams  map[string]string //路由参数缓存
	QueryParams map[string]string //url get参数缓存
	FormParams  map[string]string //表单参数缓存

	MatchedRoute string //完整的命中的路由
}
