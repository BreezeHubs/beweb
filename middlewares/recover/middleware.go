package recover

import (
	"fmt"
	"github.com/BreezeHubs/beweb"
	"runtime/debug"
)

type MiddlewareBuilder struct {
	responseStatus  int
	responseContent []byte
	logWithErr      func(err any)
	logWithContext  func(ctx *beweb.Context)
	logWithStack    func(stack string)
	logPrintSort    []string //log先后顺序记录
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (m *MiddlewareBuilder) SetPanicResponse(status int, content []byte) *MiddlewareBuilder {
	if status != 0 {
		m.responseStatus = status
	}
	if content != nil && len(content) > 0 {
		m.responseContent = content
	}
	return m
}

func (m *MiddlewareBuilder) SetLogWithErr(fn func(err any)) *MiddlewareBuilder {
	m.logWithErr = fn
	m.addOrDeleteSort("logWithErr")
	return m
}

func (m *MiddlewareBuilder) SetLogWithContext(fn func(ctx *beweb.Context)) *MiddlewareBuilder {
	m.logWithContext = fn
	m.addOrDeleteSort("logWithContext")
	return m
}

func (m *MiddlewareBuilder) SetLogWithStack(fn func(stack string)) *MiddlewareBuilder {
	m.logWithStack = fn
	m.addOrDeleteSort("logWithStack")
	return m
}

// 添加输出的排序 或 已有排序，删除原排序索，引追加新排序
func (m *MiddlewareBuilder) addOrDeleteSort(logType string) {
	for index, s := range m.logPrintSort {
		if s == logType {
			m.logPrintSort = append(m.logPrintSort[:index], m.logPrintSort[index+1:]...)
		}
	}
	m.logPrintSort = append(m.logPrintSort, logType)
}

func (m *MiddlewareBuilder) Build() beweb.Middleware {
	return func(next beweb.HandleFunc) beweb.HandleFunc {
		return func(ctx *beweb.Context) {
			defer func() {
				if err := recover(); err != nil {
					//panic后根据设置的response返回
					ctx.ResponseContent = m.responseContent
					ctx.ResponseStatus = m.responseStatus

					//按设置顺序输出
					for _, sort := range m.logPrintSort {
						switch sort {
						case "logWithErr":
							if m.logWithErr != nil {
								m.logWithErr(err)
							}
						case "logWithContext":
							if m.logWithContext != nil {
								m.logWithContext(ctx)
							}
						case "logWithStack":
							if m.logWithStack != nil {
								m.logWithStack(fmt.Sprintf("%s", debug.Stack()))
							}
						}
					}
				}
			}()

			next(ctx)
		}
	}
}
