package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Loe1210/personal-site/internal/xerrors"
	"github.com/Loe1210/personal-site/internal/xresilience"
	kitexauth "github.com/Loe1210/personal-site/kitex_gen/auth"
	"github.com/Loe1210/personal-site/kitex_gen/auth/authservice"
)

type Options struct {
	Timeout time.Duration
	Breaker *xresilience.CircuitBreaker
}

type KitexClient struct {
	cli     authservice.Client
	timeout time.Duration
	breaker *xresilience.CircuitBreaker
}

func NewKitexClient(cli authservice.Client) *KitexClient {
	return NewKitexClientWithOptions(cli, Options{Timeout: xresilience.DefaultAuthRPCTimeout})
}

func NewKitexClientWithTimeout(cli authservice.Client, timeout time.Duration) *KitexClient {
	return NewKitexClientWithOptions(cli, Options{Timeout: timeout})
}

func NewKitexClientWithOptions(cli authservice.Client, opts Options) *KitexClient {
	if opts.Timeout <= 0 {
		opts.Timeout = xresilience.DefaultAuthRPCTimeout
	}
	if opts.Breaker == nil {
		opts.Breaker = xresilience.NewCircuitBreaker(xresilience.CircuitBreakerConfig{Name: "auth-rpc"})
	}
	return &KitexClient{cli: cli, timeout: opts.Timeout, breaker: opts.Breaker}
}

func (c *KitexClient) ValidateSession(ctx context.Context, sessionID string) error {
	_, err := c.AuthContext(ctx, sessionID)
	return err
}

func (c *KitexClient) AuthContext(ctx context.Context, sessionID string) (*Context, error) {
	rpcCtx, cancel := xresilience.WithTimeout(ctx, c.timeout)
	defer cancel()
	var resp *kitexauth.AuthContext
	err := c.breaker.Run(func() error {
		return xresilience.DoWithRetry(rpcCtx, authReadRetryPolicy(), func(attemptCtx context.Context) error {
			var callErr error
			resp, callErr = c.cli.ValidateSession(attemptCtx, &kitexauth.ValidateSessionRequest{SessionId: sessionID})
			return callErr
		})
	})
	if err != nil {
		logAuthBreaker(c.breaker, err)
		if errors.Is(err, xresilience.ErrCircuitOpen) {
			return nil, xerrors.New(xerrors.CodeGatewayCircuitOpen, "auth service circuit open")
		}
		return nil, xerrors.New(xerrors.CodeAuthUpstreamFailed, "auth service unavailable")
	}
	logAuthBreaker(c.breaker, nil)
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
	err := c.breaker.Run(func() error {
		return xresilience.DoWithRetry(rpcCtx, authReadRetryPolicy(), func(attemptCtx context.Context) error {
			var callErr error
			resp, callErr = c.cli.CheckPermission(attemptCtx, &kitexauth.CheckPermissionRequest{UserId: userID, Code: code})
			return callErr
		})
	})
	if err != nil {
		logAuthBreaker(c.breaker, err)
		if errors.Is(err, xresilience.ErrCircuitOpen) {
			return false, xerrors.New(xerrors.CodeGatewayCircuitOpen, "auth service circuit open")
		}
		return false, xerrors.New(xerrors.CodeAuthUpstreamFailed, "auth service unavailable")
	}
	logAuthBreaker(c.breaker, nil)
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
	return Context{UserID: authContext.GetUserId(), Username: authContext.GetUsername(), Roles: append([]string(nil), authContext.GetRoles()...)}
}

func authReadRetryPolicy() xresilience.RetryPolicy {
	return xresilience.RetryPolicy{MaxAttempts: xresilience.DefaultReadMaxAttempts, Backoff: xresilience.DefaultRetryBackoff, Retryable: func(err error) bool { return err != nil }}
}

func logAuthBreaker(breaker *xresilience.CircuitBreaker, err error) {
	if breaker == nil {
		return
	}
	snapshot := breaker.Snapshot()
	log.Printf("component=auth_rpc breaker=%s breaker_state=%s breaker_rejected=%d breaker_failures=%d breaker_recoveries=%d err=%v", snapshot.Name, snapshot.State, snapshot.RejectedCalls, snapshot.FailureCount, snapshot.RecoveryCount, err)
}
