package xresilience

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWithTimeoutDoesNotExtendExistingDeadline(t *testing.T) {
	parent, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()

	ctx, release := WithTimeout(parent, time.Second)
	defer release()

	parentDeadline, _ := parent.Deadline()
	gotDeadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline")
	}
	if !gotDeadline.Equal(parentDeadline) {
		t.Fatalf("expected existing deadline %s, got %s", parentDeadline, gotDeadline)
	}
}

func TestRetryOnlyRetriesIdempotentMethods(t *testing.T) {
	if !IsIdempotentReadMethod("GET") || !IsIdempotentReadMethod("HEAD") {
		t.Fatal("expected GET and HEAD to be idempotent read methods")
	}
	if IsIdempotentReadMethod("POST") || IsIdempotentReadMethod("PUT") || IsIdempotentReadMethod("DELETE") {
		t.Fatal("write methods must not be treated as retryable reads")
	}
}

func TestDoWithRetryRetriesTransientReadFailureOnce(t *testing.T) {
	attempts := 0
	policy := RetryPolicy{MaxAttempts: 2, Backoff: time.Nanosecond, Retryable: func(error) bool { return true }}

	err := DoWithRetry(context.Background(), policy, func(context.Context) error {
		attempts++
		if attempts == 1 {
			return errors.New("temporary failure")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("expected retry to recover, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestFailureStatsTracksConsecutiveFailures(t *testing.T) {
	var stats FailureStats

	stats.RecordFailure()
	stats.RecordFailure()
	snapshot := stats.Snapshot()
	if snapshot.TotalFailures != 2 || snapshot.ConsecutiveFailures != 2 {
		t.Fatalf("unexpected failure snapshot: %#v", snapshot)
	}

	stats.RecordSuccess()
	snapshot = stats.Snapshot()
	if snapshot.TotalFailures != 2 || snapshot.ConsecutiveFailures != 0 {
		t.Fatalf("expected consecutive failures to reset, got %#v", snapshot)
	}
}
