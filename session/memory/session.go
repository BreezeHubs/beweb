package memory

import (
	"context"
	"errors"
	"fmt"
	"github.com/BreezeHubs/beweb/session"
	cache "github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var (
	// sentinel error. 预定义错误
	errorKeyNotFound = errors.New("key not found")
)

type Store struct {
	cache      *cache.Cache
	expiration time.Duration
	mutex      sync.RWMutex
}

func NewStore(expiration ...time.Duration) *Store {
	t := 30 * time.Minute
	if len(expiration) > 0 {
		t = expiration[0]
	}
	return &Store{cache: cache.New(t, time.Second), expiration: t}
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session := &Session{
		id:     id,
		values: sync.Map{},
	}
	s.cache.Set(id, session, s.expiration)
	return session, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	value, ok := s.cache.Get(id)
	if !ok {
		return fmt.Errorf("该 id 对应的 session 不存在， id：%s", id)
	}
	s.cache.Set(id, value, s.expiration)
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.cache.Delete(id)
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, ok := s.cache.Get(id)
	if !ok {
		return nil, fmt.Errorf("该 id 对应的 session 不存在， id：%s", id)
	}
	return value.(*Session), nil
}

type Session struct {
	id     string
	values sync.Map
}

func (s *Session) Get(ctx context.Context, key string) (any, error) {
	value, ok := s.values.Load(key)
	if !ok {
		//return nil, fmt.Errorf("%w, key %s", errorKeyNotFound, key)
		return nil, errorKeyNotFound
	}
	return value, nil
}

func (s *Session) Set(ctx context.Context, key string, value any) error {
	s.values.Store(key, value)
	return nil
}

func (s *Session) ID() string {
	return s.id
}
