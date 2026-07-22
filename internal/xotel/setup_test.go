package xotel

import (
	"context"
	"testing"
)

func TestSetupTracerProviderRequiresServiceName(t *testing.T) {
	shutdown, err := SetupTracerProvider(context.Background(), "", "")
	if err == nil {
		t.Fatal("expected missing service name error")
	}
	if shutdown != nil {
		t.Fatalf("expected nil shutdown, got %#v", shutdown)
	}
}

func TestSetupTracerProviderReturnsNoopShutdown(t *testing.T) {
	shutdown, err := SetupTracerProvider(context.Background(), "gateway", "")
	if err != nil {
		t.Fatalf("SetupTracerProvider returned error: %v", err)
	}
	if shutdown == nil {
		t.Fatal("expected shutdown function")
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown returned error: %v", err)
	}
}
