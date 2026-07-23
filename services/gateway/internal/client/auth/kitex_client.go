package auth

import (
	"context"
	"time"

	"github.com/Loe1210/personal-site/internal/xerrors"
	"github.com/Loe1210/personal-site/internal/xresilience"
	kitexauth "github.com/Loe1210/personal-site/kitex_gen/auth"
	"github.com/Loe1210/personal-site/kitex_gen/auth/authservice"
)

type KitexClient struct {
	cli     authservice.Client
	timeout time.Duration
}

func NewKitexClient(cli authservice.Client) *KitexClient {
	return NewKitexClientWithTimeout(cli, xresilience.DefaultAuthRPCTimeout)
}

func NewKitexClientWithTimeout(cli authservice.Client, timeout time.Duration) *KitexClient {
	return &KitexClient{cli: cli, timeout: timeout}
}

func (c *KitexClient) ValidateSession(ctx context.Context, sessionID string) error {
	_, err := c.AuthContext(ctx, sessionID)
	return err
}

func (c *KitexClient) AuthContext(ctx context.Context, sessionID string) (*Context, error) {
	rpcCtx, cancel := xresilience.WithTimeout(ctx, c.timeout)
	defer cancel()
	var resp *kitexauth.AuthContext
	err := xresilience.DoWithRetry(rpcCtx, authReadRetryPolicy(), func(attemptCtx context.Context) error {
		var callErr error
		resp, callErr = c.cli.ValidateSession(attemptCtx, &kitexauth.ValidateSessionRequest{SessionId: sessionID})
		return callErr
	})
	if err != nil {
		return nil, xerrors.New(xerrors.CodeAuthUpstreamFailed, "auth service unavailable")
	}
	if err := errorFromBaseResp(resp.GetBaseResp()); err != nil {
		return nil, err
	}
	authContext := contextFromPB(resp)
	return &authContext, nil
}

func (c *KitexClient) CheckPermission(ctx context.Context, userID int64, code string) (bool, error) {
	rpcCtx, cancel := xresilience.WithTimeout(ctx, c.timeout)
	defer cancel()
	var resp *kitexauth.CheckPermissionResponse
	err := xresilience.DoWithRetry(rpcCtx, authReadRetryPolicy(), func(attemptCtx context.Context) error {
		var callErr error
		resp, callErr = c.cli.CheckPermission(attemptCtx, &kitexauth.CheckPermissionRequest{UserId: userID, Code: code})
		return callErr
	})
	if err != nil {
		return false, xerrors.New(xerrors.CodeAuthUpstreamFailed, "auth service unavailable")
	}
	if err := errorFromBaseResp(resp.GetBaseResp()); err != nil {
		return false, err
	}
	return resp.GetAllowed(), nil
}

func errorFromBaseResp(resp interface {
	GetCode() int32
	GetMsg() string
}) error {
	if resp == nil || resp.GetCode() == xerrors.CodeOK {
		return nil
	}
	return xerrors.New(resp.GetCode(), resp.GetMsg())
}

func contextFromPB(authContext *kitexauth.AuthContext) Context {
	if authContext == nil {
		return Context{}
	}
	return Context{
		UserID:   authContext.GetUserId(),
		Username: authContext.GetUsername(),
		Roles:    append([]string(nil), authContext.GetRoles()...),
	}
}

func authReadRetryPolicy() xresilience.RetryPolicy {
	return xresilience.RetryPolicy{
		MaxAttempts: xresilience.DefaultReadMaxAttempts,
		Backoff:     xresilience.DefaultRetryBackoff,
		Retryable:   func(err error) bool { return err != nil },
	}
}
