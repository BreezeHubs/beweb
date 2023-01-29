package redis

import (
	"context"
	"github.com/BreezeHubs/beweb/session"
	"github.com/go-redis/redis/v8"
	"time"
)

// Store
// redis hset == map[string]map[string]string
//
//	     |         |      |
//	session id    key   value
type Store struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewStore(expiration ...time.Duration) *Store {
	t := 30 * time.Minute
	if len(expiration) > 0 {
		t = expiration[0]
	}
	return &Store{expiration: t}
}
func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	_, err := s.client.HSet(ctx, id, id, id).Result()
	if err != nil {
		return nil, err
	}

	_, err = s.client.Expire(ctx, id, s.expiration).Result()
	if err != nil {
		return nil, err
	}

	return &Session{client: s.client}, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Remove(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	//TODO implement me
	panic("implement me")
}

// Session
type Session struct {
	client redis.Cmdable
}

func (s *Session) Get(ctx context.Context, key string) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Session) Set(ctx context.Context, key string, value any) error {
	//TODO implement me
	panic("implement me")
}

func (s *Session) ID() string {
	//TODO implement me
	panic("implement me")
}
