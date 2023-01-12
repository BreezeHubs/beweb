package beweb

import (
	"net/http"
)

var _ IHTTPServer = &HTTPServer{}

type HTTPServer struct {
	*router //存储路由树

	middlewares []Middleware

	isGracefullyExit     bool   //是否开启优雅退出，默认关闭
	isGracefullyExitFunc func() //自定义的优雅退出之前的回收操作
}

// http.Handler接口 需要定义的方法
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//框架http上下文
	ctx := &Context{Req: r, Resp: w}

	//Middlewares
	//从后往前遍历，把后一个当前一个的next构建执行链条
	root := s.serve
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		root = s.middlewares[i](root)
	}

	root(ctx) //处理路由
}

// Serve 路由解析
func (s *HTTPServer) serve(ctx *Context) {
	//查找路由，执行业务逻辑
	info, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || info.n.handler == nil { //路由不存在 或 未设置handler
		//路由没有命中，404
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("NOT FOUND"))
		return
	}

	ctx.PathParams = info.pathParams //路由参数
	ctx.MatchedRoute = info.n.path
	info.n.handler(ctx) //执行对应路由的服务
}
