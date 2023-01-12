package beweb

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

/*
路由树：
- 全静态匹配：路径每一段都严格匹配
- 支持通配符匹配：*符号
- 支持参数路由：/test/:id
- 支持正则：
*/

type router struct {
	//http method => 路由树根节点
	trees map[string]*node
}

type node struct {
	path string //路由数分支

	//【静态路由】匹配的节点
	//子path到子节点的映射
	children map[string]*node

	//【通配符路由】匹配的节点
	startChild *node

	//【参数路由】匹配的节点
	paramChild *node

	//【正则路由】匹配的节点
	regExpChild *node

	//执行方法
	handler HandleFunc
}

// 路由参数
type matchInfo struct {
	n          *node
	pathParams map[string]string
}

// 创建router
func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

// childOrCreate 查找或创建对应的path
/*
优先级：
	1、route: /test/:id                => http://127.0.0.1:8080/test/1
	2、route: /test/*\/user            => http://127.0.0.1:8080/test/abc/user
	3、route: /test/Reg(^\d{4}-\d{8}$) => http://127.0.0.1:8080/test/0931-87562388

	互相不能共存，2、3能和【静态路由】共存
*/
func (n *node) childOrCreate(path string) (*node, error) {
	//【参数路由】处理
	if path[0] == ':' {
		//判断是否同时存在【通配符路由】
		if n.startChild != nil {
			return nil, errors.New("不允许同时存在【参数路由】和【通配符路由】，已存在【通配符路由】")
		}
		//判断是否同时存在【正则路由】
		if n.regExpChild != nil {
			return nil, errors.New("不允许同时存在【参数路由】和【正则路由】，已存在【正则路由】")
		}
		n.paramChild = &node{path: path}
		return n.paramChild, nil
	}

	//【通配符路由】处理
	if path == "*" {
		//判断是否同时存在【参数路由】
		if n.paramChild != nil {
			return nil, errors.New("不允许同时存在【参数路由】和【通配符路由】，已存在【参数路由】")
		}
		//判断是否同时存在【正则路由】
		if n.regExpChild != nil {
			return nil, errors.New("不允许同时存在【通配符路由】和【正则路由】，已存在【正则路由】")
		}

		//不存在则创建
		if n.startChild == nil {
			n.startChild = &node{path: path}
		}
		return n.startChild, nil
	}

	//验证是否【正则路由】
	if checkRouteRegExp(path) {
		//判断是否同时存在【参数路由】
		if n.paramChild != nil {
			return nil, errors.New("不允许同时存在【正则路由】和【参数路由】，已存在【参数路由】")
		}
		//判断是否同时存在【通配符路由】
		if n.startChild != nil {
			return nil, errors.New("不允许同时存在【正则路由】和【通配符路由】，已存在【通配符路由】")
		}

		n.regExpChild = &node{path: getRouteRegExp(path)}
		return n.regExpChild, nil
	}

	//【静态路由】
	//不存在子节点则创建
	if n.children == nil {
		n.children = map[string]*node{}
	}

	//查找对应的path
	res, ok := n.children[path]
	if !ok {
		//不存在，需要创建
		res = &node{path: path}
		n.children[path] = res
	}
	return res, nil
}

/*
addRoute：设置为私有方法，防止用户method乱传
path 加些限制：
  - 不支持空字符串
  - 必须以 / 开头
  - 不能以 / 结尾
  - 中间不能有连续的 ///
*/
func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {
	//验证
	if path == "" {
		r.routePanicPrint(method, path, errors.New("路径不支持空字符串"))
	}
	//必须以 / 开头
	if path[0] != '/' {
		r.routePanicPrint(method, path, errors.New("路径必须以 / 开头"))
	}
	//不能以 / 结尾
	if path != "/" && path[len(path)-1] == '/' {
		r.routePanicPrint(method, path, errors.New("路径不能以 / 结尾"))
	}

	//找到树第一层：method请求方式
	root, ok := r.trees[method]
	if !ok {
		//说明还没有根节点
		root = &node{path: "/"}
		r.trees[method] = root
	}

	//根节点特殊处理
	if path == "/" {
		//检查根节点重复注册
		if root.handler != nil {
			r.routePanicPrint(method, path, errors.New("路径不能重复注册"))
		}
		root.handler = handleFunc
		return
	}

	//切割path //path[1:] 去除最前面的斜杠，防止切分出空数组
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		//中间不能有连续的 ///
		if seg == "" {
			r.routePanicPrint(method, path, errors.New("路径中间不能有连续的 //"))
		}

		//递归站准位置。中途节点不存在，就要创建
		child, err := root.childOrCreate(seg)
		if err != nil {
			r.routePanicPrint(method, path, err)
		}
		root = child
	}
	if root.handler != nil {
		r.routePanicPrint(method, path, errors.New("路径不能重复注册"))
	}
	root.handler = handleFunc
}

// findRoute 路由匹配时调用
func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	//找到顶层请求方法的树节点
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	//根节点特殊处理
	if path == "/" {
		return &matchInfo{n: root}, true
	}

	//把前置和后置的 / 都去掉
	path = strings.Trim(path, "/")

	//切割path
	segs := strings.Split(path, "/")

	var pathParams map[string]string //存路由参数
	for _, seg := range segs {
		//中间不能有连续的 ///
		if seg == "" {
			r.routePanicPrint(method, path, errors.New("路径中间不能有连续的 //"))
		}

		//递归站准位置。中途节点不存在，就要创建
		child, paramFound, found := root.childOf(seg)
		if !found {
			return nil, false
		}

		//命中参数路由
		if paramFound {
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			pathParams[child.path[1:]] = seg
		}
		root = child
	}
	//表示有这个节点，但不一定有handler
	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true
	//return root, root.handler != nil
}

// childOf 优先考虑【静态路由】匹配，再考虑【参数路由】、【通配符路由】、【正则路由】匹配
// 查找对应的path
// *node：子节点，bool：是否参数路由匹配，bool：是否通配符匹配
func (n *node) childOf(path string) (*node, bool, bool) {
	//不存在【静态路由】
	if n.children == nil {
		//【参数路由】匹配
		if n.paramChild != nil {
			return n.paramChild, true, true
		}

		//【通配符路由】匹配
		if n.startChild != nil {
			return n.startChild, true, true
		}

		//【正则路由】匹配
		if n.regExpChild != nil {
			ok, err := regexp.MatchString(n.regExpChild.path, path)
			if ok && err == nil {
				return n.regExpChild, true, true
			}
		}

		return nil, false, false //匹配失败
	}

	//存在【静态路由】
	child, ok := n.children[path]
	if !ok {
		//【参数路由】匹配
		if n.paramChild != nil {
			return n.paramChild, true, true
		}

		//【通配符路由】匹配
		if n.startChild != nil {
			return n.startChild, true, true
		}

		//【正则路由】匹配
		if n.regExpChild != nil {
			ok, err := regexp.MatchString(n.regExpChild.path, path)
			if ok && err == nil {
				return n.regExpChild, true, true
			}
		}

		return nil, false, false //匹配失败
	}
	return child, false, ok
}

var (
	regFlagStart = "Reg("
	regFlagEnd   = ")"
)

// 获取【正则路由】的正则
// "Reg(^\d{4}-\d{8}$)"  =>  "^\d{4}-\d{8}$"
func getRouteRegExp(reg string) string {
	if len(reg) < 5 {
		return ""
	}
	return reg[4 : len(reg)-1]
}

// 检查是否是【正则路由】
// "Reg(^\d{4}-\d{8}$)"  =>  "Reg("、")"
func checkRouteRegExp(reg string) bool {
	if len(reg) < 5 {
		return false
	}
	return regFlagStart == reg[:4] && reg[len(reg)-1:] == regFlagEnd
}

type RoutePanicPrintFunc func(method, path string, err error)

// 路由错误告警
func (r *router) routePanicPrint(method, path string, err error) {
	panic(fmt.Sprintf("method: %s, path: %s, error: %+v", method, path, err))
}
