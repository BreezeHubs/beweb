package session

import (
	"github.com/BreezeHubs/beweb"
	"github.com/google/uuid"
)

type Manager struct {
	Propagator
	Store
	ContextSessionKey string
}

func (m *Manager) GetSession(ctx *beweb.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}
	//尝试读取 context 的 UserValues 缓存
	value, ok := ctx.UserValues[m.ContextSessionKey]
	if ok {
		return value.(Session), nil
	}

	sId, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}

	session, err := m.Get(ctx.Req.Context(), sId)
	if err != nil {
		return nil, err
	}
	ctx.UserValues[m.ContextSessionKey] = session
	return session, err
}

func (m *Manager) InitSession(ctx *beweb.Context) (Session, error) {
	id := uuid.New().String()
	session, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}

	return session, m.Inject(id, ctx.Resp) //注入到 http
}

func (m *Manager) RefreshSession(ctx *beweb.Context) error {
	session, err := m.GetSession(ctx)
	if err != nil {
		return err
	}

	return m.Refresh(ctx.Req.Context(), session.ID())
}

func (m *Manager) RemoveSession(ctx *beweb.Context) error {
	session, err := m.GetSession(ctx)
	if err != nil {
		return err
	}

	err = m.Store.Remove(ctx.Req.Context(), session.ID())
	if err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Resp) //http
}
