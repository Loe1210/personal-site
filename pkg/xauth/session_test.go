package xauth

import (
	"testing"
	"time"

	"github.com/Loe1210/personal-site/configs"
)

func TestCreateSession(t *testing.T) {
	configs.AppConfig = nil

	sessionID, err := CreateSession(7, "loe", []string{"super_admin"})
	if err != nil {
		t.Fatalf("CreateSession returned error: %v", err)
	}

	claims, err := ParseSession(sessionID)
	if err != nil {
		t.Fatalf("ParseSession returned error: %v", err)
	}

	if claims.UserID != 7 {
		t.Fatalf("expected user id 7, got %d", claims.UserID)
	}
	if claims.Username != "loe" {
		t.Fatalf("expected username loe, got %s", claims.Username)
	}
	if len(claims.Roles) != 1 || claims.Roles[0] != "super_admin" {
		t.Fatalf("unexpected roles: %#v", claims.Roles)
	}
	if claims.ExpiresAt.Before(time.Now()) {
		t.Fatalf("expected future expiry, got %s", claims.ExpiresAt)
	}
}

func TestCreateSessionBundleUsesConfiguredCookie(t *testing.T) {
	configs.AppConfig = nil
	t.Setenv("SESSION_STORE_COOKIE_NAME", "admin_session")
	t.Setenv("SESSION_STORE_PREFIX", "bundle:")
	t.Setenv("SESSION_STORE_EXPIRE_HOUR", "4")

	bundle, err := CreateSessionBundle(8, "future-admin", []string{"editor"})
	if err != nil {
		t.Fatalf("CreateSessionBundle returned error: %v", err)
	}

	if bundle.CookieName != "admin_session" {
		t.Fatalf("expected cookie name admin_session, got %q", bundle.CookieName)
	}
	if bundle.Backend != "monolith-memory" {
		t.Fatalf("expected placeholder backend monolith-memory, got %q", bundle.Backend)
	}

	claims, err := ParseSession(bundle.SessionID)
	if err != nil {
		t.Fatalf("ParseSession returned error: %v", err)
	}
	if claims.UserID != 8 {
		t.Fatalf("expected user id 8, got %d", claims.UserID)
	}
}
