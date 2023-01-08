package beweb

import (
	"strings"
)

/*
路由树：
- 全静态匹配
- 支持通配符匹配
- 支持参数路由
*/

//全静态匹配：路径每一段都严格匹配

type router struct {
	//http method => 路由树根节点
	trees map[string]*node
}

type node struct {
	path string

	//静态匹配的节点
	//子path到子节点的映射
	children map[string]*node

	//通配符匹配的节点
	startChild *node

	//路由参数匹配的节点
	paramChild *node

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

// 查找或创建对应的path
func (n *node) childOrCreate(path string) *node {
	//参数路由处理
	if path[0] == ':' {
		//判断是否同时存在通配符路由
		if n.startChild != nil {
			panic("不允许同时存在参数路由和通配符路由，已存在通配符路由[" + path + "]")
		}

		n.paramChild = &node{path: path}
		return n.paramChild
	}

	//通配符路由处理
	if path == "*" {
		//判断是否同时存在通配符路由
		if n.paramChild != nil {
			panic("不允许同时存在参数路由和通配符路由，已存在参数路由[" + path + "]")
		}

		//不存在则创建
		if n.startChild == nil {
			n.startChild = &node{path: path}
		}
		return n.startChild
	}

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
	return res
}

// childOf 优先考虑静态匹配，再考虑通配符匹配
// 查找对应的path
// *node：子节点，bool：是否参数路由匹配，bool：是否通配符匹配
func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		//参数路由匹配
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.startChild, false, n.startChild != nil //通配符匹配
	}
	child, ok := n.children[path]
	if !ok {
		//参数路由匹配
		if n.paramChild != nil {
			return n.paramChild, true, true
		}

		//通配符匹配
		return n.startChild, false, n.startChild != nil
	}
	return child, false, ok
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
		panic("路径不支持空字符串")
	}
	//必须以 / 开头
	if path[0] != '/' {
		panic("路径必须以 / 开头")
	}
	//不能以 / 结尾
	if path != "/" && path[len(path)-1] == '/' {
		panic("路径不能以 / 结尾")
	}

	//找到树
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
			panic("路径不能重复注册[/]")
		}

		root.handler = handleFunc
		return
	}

	//切割path //path[1:] 去除最前面的斜杠
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		//中间不能有连续的 ///
		if seg == "" {
			panic("路径中间不能有连续的 //")
		}

		//递归站准位置。中途节点不存在，就要创建
		child := root.childOrCreate(seg)
		root = child
	}
	if root.handler != nil {
		panic("路径不能重复注册[" + path + "]")
	}
	root.handler = handleFunc
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	//找到树
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

	var pathParams map[string]string
	for _, seg := range segs {
		//中间不能有连续的 ///
		if seg == "" {
			panic("路径中间不能有连续的 //")
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
