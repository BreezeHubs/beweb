package beweb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_AddRoute(t *testing.T) {
	//构造路由树
	//验证路由树
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail/:id",
		},
		//{
		//	method: http.MethodGet,
		//	path:   "/*",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "/*/*",
		//},
		{
			method: http.MethodGet,
			path:   "/*/abc",
		},
		//{
		//	method: http.MethodGet,
		//	path:   "/*/abc/*",
		//},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		{
			method: http.MethodGet,
			path:   "/login/:username",
		},
		{
			method: http.MethodGet,
			path:   "/*/order",
		},
	}

	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	//断言路由树和预期的一样
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"user": {
						path: "user",
						children: map[string]*node{
							"home": {
								path:    "home",
								handler: mockHandler,
							},
						},
						handler: mockHandler,
					},
					"order": {
						path: "order",
						children: map[string]*node{
							"detail": {
								path:    "detail",
								handler: mockHandler,
								paramChild: &node{
									path:    ":id",
									handler: mockHandler,
								},
							},
						},
						startChild: &node{
							path:    "*",
							handler: mockHandler,
						},
					},
					"login": {
						path: "login",
						paramChild: &node{
							path:    ":username",
							handler: mockHandler,
						},
					},
				},
				startChild: &node{
					path: "*",
					children: map[string]*node{
						"abc": {
							path:    "abc",
							handler: mockHandler,
						},
						"order": {
							path:    "order",
							handler: mockHandler,
						},
					},
				},
				handler: mockHandler,
			},
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
					"order": {
						path: "order",
						children: map[string]*node{
							"create": {
								path:    "create",
								handler: mockHandler,
							},
						},
					},
					"login": {
						path:    "login",
						handler: mockHandler,
					},
				},
			},
		},
	}

	msg, ok := wantRouter.equal(r)
	fmt.Println("ok: ", ok)
	fmt.Println("msg: ", msg)
	//assert.True(t, ok, msg)

	r = newRouter()
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	}, "路径必须以 / 开头")

	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c/", mockHandler)
	}, "路径不能以 / 结尾")

	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a///c/", mockHandler)
	}, "路径中间不能有连续的 //")

	r = newRouter()
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	}, "路径不能重复注册[/]")

	r = newRouter()
	r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	}, "路径不能重复注册[/a/b/c]")

	r = newRouter()
	r.addRoute(http.MethodGet, "/a/*", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/:id", mockHandler)
	}, "不允许同时存在参数路由和通配符路由，已存在通配符路由:id]")
}

// string错误信息
// bool代表是否相等
func (r *router) equal(y *router) (string, bool) {
	//fmt.Printf("%#v\n", y.trees["GET"].startChild)
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的http method"), false
		}
		//v, dst要相等
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点数量不相等"), false
	}

	//通配符路由处理
	if n.startChild != nil {
		msg, ok := n.startChild.equal(y.startChild)
		if !ok {
			return msg, ok
		}
	}

	if n.paramChild != nil {
		msg, ok := n.paramChild.equal(y.paramChild)
		if !ok {
			return msg, ok
		}
	}

	//比较handler
	nHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if nHandler != yHandler {
		return fmt.Sprintf("handler不相等"), false
	}

	for path, c := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点 %s 不存在", path), false
		}
		msg, ok := c.equal(dst)
		if !ok {
			return msg, false
		}
	}
	return "", true
}

func TestRouter_findRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		{
			method: http.MethodGet,
			path:   "/login/:username",
		},
	}

	r := newRouter()
	var mockHandler HandleFunc = func(ctx *Context) {}
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}
	testCases := []struct {
		name      string
		method    string
		path      string
		wantFound bool
		info      *matchInfo
	}{
		{
			//方法不存在
			name:   "method not found",
			method: http.MethodOptions,
			path:   "/order/detail",
		},
		{
			//完全命中
			name:      "order detail",
			method:    http.MethodGet,
			path:      "/order/detail",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "detail",
					handler: mockHandler,
				},
			},
		},
		{
			//命中了，但是没有handler
			name:      "order",
			method:    http.MethodGet,
			path:      "/order",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path: "order",
					children: map[string]*node{
						"detail": {
							path:    "detail",
							handler: mockHandler,
						},
					},
				},
			},
		},
		{
			name:      "order *",
			method:    http.MethodGet,
			path:      "/order/abc",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "*",
					handler: mockHandler,
				},
			},
		},
		//{
		//	//根节点
		//	name:      "root",
		//	method:    http.MethodGet,
		//	path:      "/",
		//	wantFound: true,
		//	info: &matchInfo{
		//		n: &node{
		//			path:    "/",
		//			handler: mockHandler,
		//			children: map[string]*node{
		//				"order": {
		//					path: "order",
		//					children: map[string]*node{
		//						"detail": {
		//							path:    "detail",
		//							handler: mockHandler,
		//						},
		//					},
		//					startChild: &node{
		//						path:    "*",
		//						handler: mockHandler,
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		{
			//没有path
			name:   "path not found",
			method: http.MethodGet,
			path:   "/aaaagvwegfwevw",
		},
		{
			name:      "login username",
			method:    http.MethodGet,
			path:      "/login/testname",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    ":username",
					handler: mockHandler,
				},
				pathParams: map[string]string{
					"username": "testname",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			assert.Equal(t, tc.info.pathParams, info.pathParams)
			msg, ok := tc.info.n.equal(info.n)
			assert.True(t, ok, msg)
		})
	}

}
