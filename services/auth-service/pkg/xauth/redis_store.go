package xauth

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Loe1210/personal-site/internal/xresilience"
	"github.com/gomodule/redigo/redis"
)

type RedisPool interface {
	Get() redis.Conn
}

type RedisStore struct {
	pool           RedisPool
	prefix         string
	commandTimeout time.Duration
}

func NewRedisStore(pool RedisPool, prefix string) Store {
	return NewRedisStoreWithTimeout(pool, prefix, xresilience.DefaultRedisTimeout)
}

func NewRedisStoreWithTimeout(pool RedisPool, prefix string, commandTimeout time.Duration) Store {
	if prefix == "" {
		prefix = "session:"
	}
	if commandTimeout <= 0 {
		commandTimeout = xresilience.DefaultRedisTimeout
	}
	return &RedisStore{pool: pool, prefix: prefix, commandTimeout: commandTimeout}
}

func (s *RedisStore) Save(ctx context.Context, sessionID string, claims *Claims, ttl time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if s.pool == nil {
		return errors.New("redis pool is required")
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		return err
	}
	conn := s.pool.Get()
	defer conn.Close()
	seconds := int(ttl.Seconds())
	if seconds <= 0 {
		seconds = int(time.Hour.Seconds())
	}
	_, err = redis.DoWithTimeout(conn, s.commandTimeout, "SETEX", s.key(sessionID), seconds, payload)
	return err
}

func (s *RedisStore) Get(ctx context.Context, sessionID string) (*Claims, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if s.pool == nil {
		return nil, errors.New("redis pool is required")
	}
	conn := s.pool.Get()
	defer conn.Close()
	payload, err := redis.Bytes(redis.DoWithTimeout(conn, s.commandTimeout, "GET", s.key(sessionID)))
	if errors.Is(err, redis.ErrNil) {
		return nil, errSessionNotFound
	}
	if err != nil {
		return nil, err
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}
	if time.Now().After(claims.ExpiresAt) {
		_ = s.Delete(ctx, sessionID)
		return nil, errSessionNotFound
	}
	claims.Roles = append([]string(nil), claims.Roles...)
	return &claims, nil
}

func (s *RedisStore) Delete(ctx context.Context, sessionID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if s.pool == nil {
		return errors.New("redis pool is required")
	}
	conn := s.pool.Get()
	defer conn.Close()
	_, err := redis.DoWithTimeout(conn, s.commandTimeout, "DEL", s.key(sessionID))
	return err
}

func (s *RedisStore) Backend() string {
	return redisBackend
}

func (s *RedisStore) key(sessionID string) string {
	return s.prefix + sessionID
}
