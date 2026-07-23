package service

import (
	"context"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

type timeoutCapturingRedisConn struct {
	redis.Conn
	timeout time.Duration
}

func (c *timeoutCapturingRedisConn) DoWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (interface{}, error) {
	c.timeout = timeout
	return nil, redis.ErrNil
}

func (c *timeoutCapturingRedisConn) ReceiveWithTimeout(timeout time.Duration) (interface{}, error) {
	c.timeout = timeout
	return nil, redis.ErrNil
}

func (c *timeoutCapturingRedisConn) Close() error { return nil }

type timeoutCapturingRedisPool struct {
	conn *timeoutCapturingRedisConn
}

func (p *timeoutCapturingRedisPool) Get() redis.Conn { return p.conn }

func TestRedisArticleCacheGetUsesCommandTimeout(t *testing.T) {
	conn := &timeoutCapturingRedisConn{}
	cache := NewRedisArticleCacheWithTimeout(&timeoutCapturingRedisPool{conn: conn}, "article:", time.Minute, 90*time.Millisecond)

	_, _, _ = cache.GetArticle(context.Background(), 12)

	if conn.timeout != 90*time.Millisecond {
		t.Fatalf("expected Redis command timeout, got %s", conn.timeout)
	}
}
