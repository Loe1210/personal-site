package auth

import (
	"context"
	"testing"

	kitexauth "github.com/Loe1210/personal-site/kitex_gen/auth"
)

func TestContextFromPBMapsFields(t *testing.T) {
	got := contextFromPB(&kitexauth.AuthContext{
		UserId:   42,
		Username: "admin",
		Roles:    []string{"admin", "editor"},
	})

	if got.UserID != 42 || got.Username != "admin" {
		t.Fatalf("unexpected auth context: %#v", got)
	}
	if len(got.Roles) != 2 || got.Roles[0] != "admin" || got.Roles[1] != "editor" {
		t.Fatalf("unexpected roles: %#v", got.Roles)
	}
}

func TestContextFromPBNilIsZeroContext(t *testing.T) {
	got := contextFromPB(nil)
	if got.UserID != 0 || got.Username != "" || len(got.Roles) != 0 {
		t.Fatalf("expected zero context, got %#v", got)
	}
}

func TestKitexClientImplementsSessionValidator(t *testing.T) {
	var _ interface {
		ValidateSession(context.Context, string) error
	} = (*KitexClient)(nil)
}
