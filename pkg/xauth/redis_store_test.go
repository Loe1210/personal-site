package xauth

import (
	"context"
	"testing"
	"time"
)

func TestSessionUsesConfiguredStore(t *testing.T) {
	store := newMemoryStoreForTest()
	UseStore(store)
	t.Cleanup(func() { UseStore(newMemoryStoreForTest()) })

	sessionID, err := CreateSessionWithContext(context.Background(), 7, "admin", []string{"super_admin"})
	if err != nil {
		t.Fatalf("CreateSessionWithContext returned error: %v", err)
	}

	claims, err := ParseSessionWithContext(context.Background(), sessionID)
	if err != nil {
		t.Fatalf("ParseSessionWithContext returned error: %v", err)
	}
	if claims.UserID != 7 {
		t.Fatalf("expected user id 7, got %d", claims.UserID)
	}
	if claims.ExpiresAt.Before(time.Now()) {
		t.Fatal("expected future expiration")
	}
}
