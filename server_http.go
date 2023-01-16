package beweb

import (
	"net/http"
	"time"
)

var _ IHTTPServer = &HTTPServer{}

type HTTPServer struct {
	//路由组缓存
	groupNameCache        string
	groupMiddlewaresCache []Middleware

	*router                  //存储路由树
	middlewares []Middleware //公共middlewares

	//config
	isGracefullyExit      bool          //是否开启优雅退出，默认关闭，默认false
	isGracefullyExitFunc  func()        //自定义的优雅退出之前的回收操作，默认nil
	gracefullyExitTimeout time.Duration //优雅退出超时，默认30s
	shutdownTimeout       time.Duration //http退出超时，默认10s

	templateEngine TemplateEngine
}

// http.Handler接口 需要定义的方法
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//框架http上下文
	ctx := &Context{Req: r, Resp: w, templateEngine: s.templateEngine}

	//Middlewares
	//从后往前遍历，把后一个当前一个的next构建执行链条
	root := s.serve
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		root = s.middlewares[i](root)
	}

	//将response刷新给ctx.Resp
	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			s.flashResponseData(ctx) //将response刷新给ctx.Resp
		}
	}

	root = m(root)
	root(ctx) //处理路由
}

// Serve 路由解析
func (s *HTTPServer) serve(ctx *Context) {
	//查找路由，执行业务逻辑
	info, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || info.n.handler == nil { //路由不存在 或 未设置handler
		//路由没有命中，404
		ctx.ResponseStatus = 404
		ctx.ResponseContent = []byte("NOT FOUND")
		return
	}

	//初始化body、Form、url get、请求头参数
	ctx.initBody()
	ctx.initFormValue()
	ctx.initQueryValue()
	ctx.initHeaderValue()

	ctx.PathParams = info.pathParams //路由参数
	ctx.MatchedRoute = info.n.route  //完整的命中的路由

	//Middlewares
	//从后往前遍历，把后一个当前一个的next构建执行链条
	root := info.n.handler
	mdls := info.n.middlewares
	if ok {
		for i := len(mdls) - 1; i >= 0; i-- {
			root = mdls[i](root)
		}
	}
	root(ctx) //执行对应路由的服务
}

// 将response刷新给ctx.Resp
func (s *HTTPServer) flashResponseData(ctx *Context) {
	if ctx.ResponseStatus != 0 {
		ctx.Resp.WriteHeader(ctx.ResponseStatus)
	}
	if ctx.ResponseContent != nil && len(ctx.ResponseContent) > 0 {
		_, _ = ctx.Resp.Write(ctx.ResponseContent)
	}
}
