package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Loe1210/personal-site/services/auth-service/internal/model"
	"github.com/Loe1210/personal-site/services/auth-service/pkg/xauth"
)

type fakeUserRepository struct {
	user        *model.User
	roles       []string
	permissions map[string]bool
}

func (f *fakeUserRepository) Login(_ context.Context, username, password string) (*model.User, []string, error) {
	if username != "admin" || password != "admin" {
		return nil, nil, errors.New("invalid credentials")
	}
	return f.user, f.roles, nil
}

func (f *fakeUserRepository) GetByID(_ context.Context, userID int64) (*model.User, error) {
	if f.user == nil || f.user.ID != userID {
		return nil, errors.New("user not found")
	}
	return f.user, nil
}

func (f *fakeUserRepository) HasPermission(_ context.Context, userID int64, code string) (bool, error) {
	if f.user == nil || f.user.ID != userID {
		return false, errors.New("user not found")
	}
	return f.permissions[code], nil
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		user:  &model.User{ID: 1, Username: "admin", Nickname: "Admin"},
		roles: []string{"super_admin"},
		permissions: map[string]bool{
			"article:read": true,
		},
	}
}

func TestCreateSessionFromCredentials(t *testing.T) {
	svc := NewAuthService(newFakeUserRepository())

	resp, err := svc.CreateSession(context.Background(), "admin", "admin")
	if err != nil {
		t.Fatalf("CreateSession returned error: %v", err)
	}
	if resp.SessionID == "" {
		t.Fatal("expected session id")
	}
	if resp.Username != "admin" {
		t.Fatalf("expected username admin, got %s", resp.Username)
	}
}

func TestValidateSessionReturnsAuthContext(t *testing.T) {
	svc := NewAuthService(newFakeUserRepository())
	bundle, err := xauth.CreateSessionBundle(1, "admin", []string{"super_admin"})
	if err != nil {
		t.Fatalf("CreateSessionBundle returned error: %v", err)
	}

	ctx, err := svc.ValidateSession(context.Background(), bundle.SessionID)
	if err != nil {
		t.Fatalf("ValidateSession returned error: %v", err)
	}
	if ctx.UserID != 1 || len(ctx.Roles) != 1 || ctx.Roles[0] != "super_admin" {
		t.Fatalf("unexpected auth context: %#v", ctx)
	}
}

func TestCheckPermissionDelegatesToRepository(t *testing.T) {
	svc := NewAuthService(newFakeUserRepository())

	allowed, err := svc.CheckPermission(context.Background(), 1, "article:read")
	if err != nil {
		t.Fatalf("CheckPermission returned error: %v", err)
	}
	if !allowed {
		t.Fatal("expected article:read to be allowed")
	}
}
