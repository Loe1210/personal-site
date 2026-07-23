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

func TestNewResolverReturnsNilForEmptyAddr(t *testing.T) {
	resolver, err := NewResolver("")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resolver != nil {
		t.Fatalf("expected nil resolver for empty addr, got %#v", resolver)
	}
}

func TestNewResolverRejectsInvalidAddr(t *testing.T) {
	resolver, err := NewResolver("not-a-host-port")
	if err == nil {
		t.Fatal("expected invalid addr error")
	}
	if resolver != nil {
		t.Fatalf("expected nil resolver, got %#v", resolver)
	}
}
