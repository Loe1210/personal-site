package xsafe

import (
	"errors"
	"testing"
	"time"
)

func TestGoCapturesPanic(t *testing.T) {
	errCh := make(chan error, 1)
	SetPanicHandler(func(value any) {
		errCh <- errors.New("captured")
	})
	defer SetPanicHandler(nil)

	Go(func() {
		panic("boom")
	})

	select {
	case err := <-errCh:
		if err.Error() != "captured" {
			t.Fatalf("unexpected panic handler error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("expected panic handler to be called")
	}
}

func TestDeferRecoverCapturesPanic(t *testing.T) {
	var captured any
	func() {
		defer DeferRecover(func(value any) {
			captured = value
		})
		panic("boom")
	}()

	if captured != "boom" {
		t.Fatalf("expected panic value boom, got %#v", captured)
	}
}
