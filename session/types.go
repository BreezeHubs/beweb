package session

import (
	"context"
	"net/http"
)

// Store 管理session本身、
type Store interface {
	// Generate
	// sesssion id、timeout可以交予用户决定
	Generate(ctx context.Context, id string) (Session, error)

	Refresh(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (Session, error)
}

type Session interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any) error
	ID() string
}

type Propagator interface {
	// Inject 将 session id 注入到 http
	Inject(id string, writer http.ResponseWriter) error

	// Extract 将 session id 从 http 提取出来
	Extract(req *http.Request) (string, error)

	// Remove 将 session id 从 http 删除
	Remove(writer http.ResponseWriter) error
}
