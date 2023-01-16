package beweb

import (
	"errors"
	"fmt"
	"strings"
)

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
func (r *router) routePanicPrint(group, method, path string, err error) {
	panic(fmt.Sprintf("group: %s, method: %s, path: %s, error: %+v", group, method, path, err))
}

// 检查group格式
func (r *router) checkGroupName(group, method, path string) {
	if group != "" {
		if group[0] != '/' {
			r.routePanicPrint(group, method, path, errors.New("group 必须以 / 开头"))
		}
		if group != "/" && group[len(group)-1] == '/' {
			r.routePanicPrint(group, method, path, errors.New("group 不能以 / 结尾"))
		}
		segs := strings.Split(path[1:], "/")
		for _, seg := range segs {
			//中间不能有连续的 ///
			if seg == "" {
				r.routePanicPrint(group, method, path, errors.New("group 中间不能有连续的 //"))
			}
		}
		path = group + path
	}
}
