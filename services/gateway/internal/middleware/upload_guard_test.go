package middleware

import (
	"errors"
	"testing"
)

func TestUploadGuardRejectsTooLargeRequest(t *testing.T) {
	guard := NewUploadGuard(UploadGuardConfig{MaxBodyBytes: 4, MaxConcurrent: 1})
	_, err := guard.Acquire(5)
	if !errors.Is(err, ErrUploadTooLarge) {
		t.Fatalf("expected upload too large error, got %v", err)
	}
}

func TestUploadGuardLimitsConcurrentUploads(t *testing.T) {
	guard := NewUploadGuard(UploadGuardConfig{MaxBodyBytes: 100, MaxConcurrent: 1})
	release, err := guard.Acquire(10)
	if err != nil {
		t.Fatalf("first acquire should pass: %v", err)
	}
	defer release()

	_, err = guard.Acquire(10)
	if !errors.Is(err, ErrUploadBusy) {
		t.Fatalf("expected concurrent upload limit error, got %v", err)
	}
}
