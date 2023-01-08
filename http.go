package beweb

import "net/http"

var _ IHTTPServer = &HTTPServer{}

type HTTPServer struct {
	*router //存储路由树
}

// http.Handler接口 需要定义的方法
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//框架http上下文
	ctx := &Context{Req: r, Resp: w}
	s.Serve(ctx) //处理路由
}

// Serve 路由解析
func (s *HTTPServer) Serve(ctx *Context) {
	//查找路由，执行业务逻辑
	info, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || info.n.handler == nil { //路由不存在 或 未设置handler
		//路由没有命中，404
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("NOT FOUND"))
		return
	}

	ctx.PathParams = info.pathParams //路由参数
	info.n.handler(ctx)              //执行对应路由的服务
}
