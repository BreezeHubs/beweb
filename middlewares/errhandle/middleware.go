package errhandle

import "github.com/BreezeHubs/beweb"

type MiddlewareBuilder struct {
	responseContent map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{responseContent: make(map[int][]byte, 16)}
}

// AddResponseContent 增加响应
func (m *MiddlewareBuilder) AddResponseContent(status int, content []byte) *MiddlewareBuilder {
	m.responseContent[status] = content
	return m
}

func (m *MiddlewareBuilder) Build() beweb.Middleware {
	return func(next beweb.HandleFunc) beweb.HandleFunc {
		return func(ctx *beweb.Context) {
			next(ctx)

			resp, ok := m.responseContent[ctx.ResponseStatus]
			if ok {
				ctx.ResponseContent = resp
			}
		}
	}
}
