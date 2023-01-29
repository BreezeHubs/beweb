package session

import (
	"github.com/BreezeHubs/beweb"
	"net/http"
	"testing"
)

func TestSession(t *testing.T) {
	var (
		p Propagator
		s Store
		m *Manager
	)

	server := beweb.NewHTTPServer(
		//登录校验
		beweb.WithMiddlewares(func(next beweb.HandleFunc) beweb.HandleFunc {
			return func(ctx *beweb.Context) {
				if ctx.Req.URL.Path == "/login" {
					next(ctx)
					return
				}

				sId, err := p.Extract(ctx.Req)
				if err != nil {
					ctx.Response(http.StatusUnauthorized, []byte("请重新登录"))
					return
				}

				_, err = s.Get(ctx.Req.Context(), sId)
				if err != nil {
					ctx.Response(http.StatusUnauthorized, []byte("请重新登录"))
					return
				}

				//刷新 session 的过期时间
				err = m.RefreshSession(ctx)
				if err != nil {
					ctx.Response(http.StatusUnauthorized, []byte("请重新登录"))
					return
				}

				next(ctx)
			}
		}),
	)

	server.Get("login", func(ctx *beweb.Context) {
		session, err := m.InitSession(ctx)
		if err != nil {
			ctx.Response(http.StatusInternalServerError, []byte("登录异常"+err.Error()))
			return
		}

		err = session.Set(ctx.Req.Context(), "nickname", "breeze")
		if err != nil {
			ctx.Response(http.StatusInternalServerError, []byte("登录异常"+err.Error()))
			return
		}

		ctx.Response(http.StatusOK, []byte("登录成功"))
	})

	server.Get("/logout", func(ctx *beweb.Context) {
		err := m.RemoveSession(ctx)
		if err != nil {
			ctx.Response(http.StatusInternalServerError, []byte("退出登录异常"+err.Error()))
			return
		}

		ctx.Response(http.StatusOK, []byte("退出登录成功"))
	})

	server.Get("/manager", func(ctx *beweb.Context) {
		session, err := m.GetSession(ctx)
		if err != nil {
			ctx.Response(http.StatusUnauthorized, []byte("请重新登录"+err.Error()))
			return
		}

		nickname, err := session.Get(ctx.Req.Context(), "nickname")
		if err != nil {
			ctx.Response(http.StatusInternalServerError, []byte("登录信息异常："+err.Error()))
			return
		}
		ctx.Response(http.StatusOK, []byte("nickname："+nickname.(string)))
	})

	server.Start(":8080")
}
