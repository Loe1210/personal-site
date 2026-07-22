package rpc

import (
	"context"
	"testing"

	"github.com/Loe1210/personal-site/internal/xerrors"
	kitexauth "github.com/Loe1210/personal-site/kitex_gen/auth"
	kitexbase "github.com/Loe1210/personal-site/kitex_gen/base"
	"github.com/Loe1210/personal-site/services/auth-service/internal/model"
	"github.com/Loe1210/personal-site/services/auth-service/internal/service"
)

type fakeUserRepository struct{}

func (f *fakeUserRepository) Login(context.Context, string, string) (*model.User, []string, error) {
	return &model.User{ID: 1, Username: "admin"}, []string{"admin"}, nil
}

func (f *fakeUserRepository) GetByID(context.Context, int64) (*model.User, error) {
	return &model.User{ID: 1, Username: "admin"}, nil
}

func (f *fakeUserRepository) HasPermission(context.Context, int64, string) (bool, error) {
	return false, nil
}

func TestValidateSessionBusinessErrorReturnsBaseResp(t *testing.T) {
	handler := NewHandler(service.NewAuthService(&fakeUserRepository{}))

	resp, err := handler.ValidateSession(context.Background(), &kitexauth.ValidateSessionRequest{SessionId: "expired"})

	if err != nil {
		t.Fatalf("expected nil rpc error, got %v", err)
	}
	if resp.GetBaseResp().GetCode() != xerrors.CodeAuthSessionExpired {
		t.Fatalf("expected expired session code, got %#v", resp.GetBaseResp())
	}
	if resp.GetBaseResp().GetMsg() == "" {
		t.Fatalf("expected business message, got %#v", resp.GetBaseResp())
	}
}

func TestCheckPermissionDeniedReturnsBaseResp(t *testing.T) {
	handler := NewHandler(service.NewAuthService(&fakeUserRepository{}))

	resp, err := handler.CheckPermission(context.Background(), &kitexauth.CheckPermissionRequest{UserId: 1, Code: "article:write"})

	if err != nil {
		t.Fatalf("expected nil rpc error, got %v", err)
	}
	if resp.GetBaseResp().GetCode() != xerrors.CodeAuthPermissionDenied {
		t.Fatalf("expected permission denied code, got %#v", resp.GetBaseResp())
	}
	if resp.GetAllowed() {
		t.Fatal("expected denied permission")
	}
}

func TestAuthContextSuccessBaseRespIsOK(t *testing.T) {
	resp := &kitexauth.AuthContext{BaseResp: &kitexbase.BaseResp{Code: xerrors.CodeOK, Msg: "success"}}
	if resp.GetBaseResp().GetCode() != xerrors.CodeOK {
		t.Fatalf("expected success base resp, got %#v", resp.GetBaseResp())
	}
}
