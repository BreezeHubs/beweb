package beweb

import (
	"net/http"
)

type Context struct {
	Req         *http.Request
	Resp        http.ResponseWriter
	PathParams  map[string]string
	QueryParams map[string]string
	FormParams  map[string]string
}
