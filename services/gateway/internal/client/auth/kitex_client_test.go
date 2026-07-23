package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	kitexauth "github.com/Loe1210/personal-site/kitex_gen/auth"
	kitexbase "github.com/Loe1210/personal-site/kitex_gen/base"
	"github.com/cloudwego/kitex/client/callopt"

	"github.com/Loe1210/personal-site/internal/xerrors"
	"github.com/Loe1210/personal-site/internal/xresilience"
)

type fakeAuthServiceClient struct {
	validateResp  *kitexauth.AuthContext
	validateErr   error
	validateCalls int
}

func (f *fakeAuthServiceClient) ValidateSession(ctx context.Context, req *kitexauth.ValidateSessionRequest, callOptions ...callopt.Option) (*kitexauth.AuthContext, error) {
	f.validateCalls++
	return f.validateResp, f.validateErr
}

func (f *fakeAuthServiceClient) CheckPermission(ctx context.Context, req *kitexauth.CheckPermissionRequest, callOptions ...callopt.Option) (*kitexauth.CheckPermissionResponse, error) {
	return &kitexauth.CheckPermissionResponse{BaseResp: &kitexbase.BaseResp{Code: xerrors.CodeOK, Msg: "success"}, Allowed: true}, nil
}

func TestContextFromPBMapsFields(t *testing.T) {
	got := contextFromPB(&kitexauth.AuthContext{UserId: 42, Username: "admin", Roles: []string{"admin", "editor"}})

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

func TestAuthContextMapsBaseRespBusinessError(t *testing.T) {
	client := NewKitexClient(&fakeAuthServiceClient{validateResp: &kitexauth.AuthContext{BaseResp: &kitexbase.BaseResp{Code: xerrors.CodeAuthSessionExpired, Msg: "login expired"}}})

	_, err := client.AuthContext(context.Background(), "expired")

	if xerrors.CodeOf(err) != xerrors.CodeAuthSessionExpired {
		t.Fatalf("expected expired session app error, got %v", err)
	}
}

func TestAuthContextMapsTransportErrorToUpstreamError(t *testing.T) {
	client := NewKitexClient(&fakeAuthServiceClient{validateErr: errors.New("dial failed")})

	_, err := client.AuthContext(context.Background(), "session")

	if xerrors.CodeOf(err) != xerrors.CodeAuthUpstreamFailed {
		t.Fatalf("expected upstream failed app error, got %v", err)
	}
}

func TestAuthContextOpenBreakerSkipsRPCAndReturnsUnifiedError(t *testing.T) {
	inner := &fakeAuthServiceClient{validateResp: &kitexauth.AuthContext{BaseResp: &kitexbase.BaseResp{Code: xerrors.CodeOK, Msg: "success"}}}
	breaker := xresilience.NewCircuitBreaker(xresilience.CircuitBreakerConfig{Name: "auth-rpc", FailureThreshold: 1, OpenTimeout: time.Minute})
	breaker.RecordFailure()
	client := NewKitexClientWithOptions(inner, Options{Timeout: time.Second, Breaker: breaker})

	_, err := client.AuthContext(context.Background(), "session")

	if xerrors.CodeOf(err) != xerrors.CodeGatewayCircuitOpen {
		t.Fatalf("expected circuit open app error, got %v", err)
	}
	if inner.validateCalls != 0 {
		t.Fatalf("expected open breaker to skip RPC, got %d calls", inner.validateCalls)
	}
	if breaker.Snapshot().RejectedCalls != 1 {
		t.Fatalf("expected rejected call to be counted, got %#v", breaker.Snapshot())
	}
}
