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

var ErrCircuitOpen = errors.New("circuit breaker is open")

type CancelFunc = context.CancelFunc

type BreakerState string

const (
	BreakerClosed   BreakerState = "closed"
	BreakerOpen     BreakerState = "open"
	BreakerHalfOpen BreakerState = "half_open"
)

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

type CircuitBreakerConfig struct {
	Name             string
	FailureThreshold int64
	OpenTimeout      time.Duration
	Now              func() time.Time
}

type CircuitBreakerSnapshot struct {
	Name                string
	State               BreakerState
	FailureCount        int64
	RejectedCalls       int64
	RecoveryCount       int64
	LastFailureAt       time.Time
	LastStateTransition time.Time
}

type CircuitBreaker struct {
	mu                  sync.Mutex
	name                string
	failureThreshold    int64
	openTimeout         time.Duration
	now                 func() time.Time
	state               BreakerState
	failureCount        int64
	rejectedCalls       int64
	recoveryCount       int64
	lastFailureAt       time.Time
	lastStateTransition time.Time
}

func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 3
	}
	if cfg.OpenTimeout <= 0 {
		cfg.OpenTimeout = 5 * time.Second
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	return &CircuitBreaker{name: cfg.Name, failureThreshold: cfg.FailureThreshold, openTimeout: cfg.OpenTimeout, now: cfg.Now, state: BreakerClosed, lastStateTransition: cfg.Now()}
}

func (b *CircuitBreaker) Allow() bool {
	if b == nil {
		return true
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == BreakerOpen && !b.now().Before(b.lastStateTransition.Add(b.openTimeout)) {
		b.state = BreakerHalfOpen
		b.lastStateTransition = b.now()
	}
	if b.state == BreakerOpen {
		b.rejectedCalls++
		return false
	}
	return true
}

func (b *CircuitBreaker) Run(operation func() error) error {
	if b == nil {
		return operation()
	}
	if !b.Allow() {
		return ErrCircuitOpen
	}
	err := operation()
	if err != nil {
		b.RecordFailure()
		return err
	}
	b.RecordSuccess()
	return nil
}

func (b *CircuitBreaker) RecordFailure() {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	now := b.now()
	b.failureCount++
	b.lastFailureAt = now
	if b.state == BreakerHalfOpen || b.failureCount >= b.failureThreshold {
		b.state = BreakerOpen
		b.lastStateTransition = now
	}
}

func (b *CircuitBreaker) RecordSuccess() {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == BreakerHalfOpen {
		b.recoveryCount++
	}
	b.failureCount = 0
	if b.state != BreakerClosed {
		b.state = BreakerClosed
		b.lastStateTransition = b.now()
	}
}

func (b *CircuitBreaker) Snapshot() CircuitBreakerSnapshot {
	if b == nil {
		return CircuitBreakerSnapshot{State: BreakerClosed}
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return CircuitBreakerSnapshot{Name: b.name, State: b.state, FailureCount: b.failureCount, RejectedCalls: b.rejectedCalls, RecoveryCount: b.recoveryCount, LastFailureAt: b.lastFailureAt, LastStateTransition: b.lastStateTransition}
}
