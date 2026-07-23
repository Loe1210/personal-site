package xresilience

import (
	"context"
	"errors"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	DefaultGatewayProxyTimeout = 6 * time.Second
	DefaultAuthRPCTimeout      = 800 * time.Millisecond
	DefaultRedisTimeout        = 300 * time.Millisecond
	DefaultMySQLTimeout        = 1500 * time.Millisecond
	DefaultRetryBackoff        = 50 * time.Millisecond
	DefaultReadMaxAttempts     = 2
)

type CancelFunc = context.CancelFunc

func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, CancelFunc) {
	if timeout <= 0 {
		return context.WithCancel(ctx)
	}
	if deadline, ok := ctx.Deadline(); ok && time.Until(deadline) <= timeout {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, timeout)
}

func IsIdempotentReadMethod(method string) bool {
	switch strings.ToUpper(method) {
	case "GET", "HEAD":
		return true
	default:
		return false
	}
}

type RetryPolicy struct {
	MaxAttempts int
	Backoff     time.Duration
	Retryable   func(error) bool
}

func (p RetryPolicy) normalized() RetryPolicy {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	if p.Backoff < 0 {
		p.Backoff = 0
	}
	if p.Retryable == nil {
		p.Retryable = func(error) bool { return false }
	}
	return p
}

func DoWithRetry(ctx context.Context, policy RetryPolicy, operation func(context.Context) error) error {
	policy = policy.normalized()
	var lastErr error
	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = operation(ctx)
		if lastErr == nil {
			return nil
		}
		if attempt == policy.MaxAttempts || !policy.Retryable(lastErr) {
			return lastErr
		}
		if policy.Backoff <= 0 {
			continue
		}
		timer := time.NewTimer(policy.Backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return lastErr
}

func IsTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, os.ErrDeadlineExceeded) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "timeout")
}

type FailureSnapshot struct {
	TotalFailures       int64
	ConsecutiveFailures int64
	LastFailureAt       time.Time
}

type FailureStats struct {
	mu                  sync.Mutex
	totalFailures       int64
	consecutiveFailures int64
	lastFailureAt       time.Time
}

func (s *FailureStats) RecordFailure() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.totalFailures++
	s.consecutiveFailures++
	s.lastFailureAt = time.Now()
}

func (s *FailureStats) RecordSuccess() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.consecutiveFailures = 0
}

func (s *FailureStats) Snapshot() FailureSnapshot {
	s.mu.Lock()
	defer s.mu.Unlock()
	return FailureSnapshot{TotalFailures: s.totalFailures, ConsecutiveFailures: s.consecutiveFailures, LastFailureAt: s.lastFailureAt}
}
