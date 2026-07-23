package auth

import (
	"context"
	"testing"
	"time"

	kitexauth "github.com/Loe1210/personal-site/kitex_gen/auth"
	kitexbase "github.com/Loe1210/personal-site/kitex_gen/base"
	"github.com/cloudwego/kitex/client/callopt"

	"github.com/Loe1210/personal-site/internal/xerrors"
)

type deadlineCapturingAuthClient struct {
	deadlineSet bool
}

func (c *deadlineCapturingAuthClient) ValidateSession(ctx context.Context, req *kitexauth.ValidateSessionRequest, callOptions ...callopt.Option) (*kitexauth.AuthContext, error) {
	_, c.deadlineSet = ctx.Deadline()
	return &kitexauth.AuthContext{BaseResp: &kitexbase.BaseResp{Code: xerrors.CodeOK, Msg: "success"}}, nil
}

func (c *deadlineCapturingAuthClient) CheckPermission(ctx context.Context, req *kitexauth.CheckPermissionRequest, callOptions ...callopt.Option) (*kitexauth.CheckPermissionResponse, error) {
	_, c.deadlineSet = ctx.Deadline()
	return &kitexauth.CheckPermissionResponse{BaseResp: &kitexbase.BaseResp{Code: xerrors.CodeOK, Msg: "success"}, Allowed: true}, nil
}

func TestAuthContextAddsRPCDeadline(t *testing.T) {
	inner := &deadlineCapturingAuthClient{}
	client := NewKitexClientWithTimeout(inner, 50*time.Millisecond)

	_, err := client.AuthContext(context.Background(), "session")

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if !inner.deadlineSet {
		t.Fatal("expected auth RPC context to have a deadline")
	}
}
