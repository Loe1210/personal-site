package xnacos

import (
	"context"
	"testing"
)

func TestRegisterServiceRequiresServiceName(t *testing.T) {
	client := NewMemoryClient()
	err := client.RegisterService(context.Background(), "", "127.0.0.1", 9001)
	if err == nil {
		t.Fatal("expected error")
	}
}
