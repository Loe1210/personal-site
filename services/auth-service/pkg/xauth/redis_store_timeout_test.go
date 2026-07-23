package xauth

import (
	"context"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

type timeoutCapturingConn struct {
	redis.Conn
	timeout time.Duration
}

func (c *timeoutCapturingConn) DoWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (interface{}, error) {
	c.timeout = timeout
	return nil, redis.ErrNil
}

func (c *timeoutCapturingConn) ReceiveWithTimeout(timeout time.Duration) (interface{}, error) {
	c.timeout = timeout
	return nil, redis.ErrNil
}

func (c *timeoutCapturingConn) Close() error { return nil }

type timeoutCapturingPool struct {
	conn *timeoutCapturingConn
}

func (p *timeoutCapturingPool) Get() redis.Conn { return p.conn }

func TestRedisStoreGetUsesCommandTimeout(t *testing.T) {
	conn := &timeoutCapturingConn{}
	store := NewRedisStoreWithTimeout(&timeoutCapturingPool{conn: conn}, "session:", 75*time.Millisecond)

	_, _ = store.Get(context.Background(), "missing")

	if conn.timeout != 75*time.Millisecond {
		t.Fatalf("expected Redis command timeout, got %s", conn.timeout)
	}
}
