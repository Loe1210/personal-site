package xnacos

import "testing"

func TestNewRegistryReturnsNilForEmptyAddr(t *testing.T) {
	registry, err := NewRegistry("")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if registry != nil {
		t.Fatalf("expected nil registry for empty addr, got %#v", registry)
	}
}

func TestNewRegistryRejectsInvalidAddr(t *testing.T) {
	registry, err := NewRegistry("not-a-host-port")
	if err == nil {
		t.Fatal("expected invalid addr error")
	}
	if registry != nil {
		t.Fatalf("expected nil registry, got %#v", registry)
	}
}
