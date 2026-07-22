package auth

import (
	"context"

	kitexauth "github.com/Loe1210/personal-site/kitex_gen/auth"
	"github.com/Loe1210/personal-site/kitex_gen/auth/authservice"
)

type KitexClient struct {
	cli authservice.Client
}

func NewKitexClient(cli authservice.Client) *KitexClient {
	return &KitexClient{cli: cli}
}

func (c *KitexClient) ValidateSession(ctx context.Context, sessionID string) error {
	_, err := c.AuthContext(ctx, sessionID)
	return err
}

func (c *KitexClient) AuthContext(ctx context.Context, sessionID string) (*Context, error) {
	resp, err := c.cli.ValidateSession(ctx, &kitexauth.ValidateSessionRequest{SessionId: sessionID})
	if err != nil {
		return nil, err
	}
	authContext := contextFromPB(resp)
	return &authContext, nil
}

func (c *KitexClient) CheckPermission(ctx context.Context, userID int64, code string) (bool, error) {
	resp, err := c.cli.CheckPermission(ctx, &kitexauth.CheckPermissionRequest{UserId: userID, Code: code})
	if err != nil {
		return false, err
	}
	return resp.GetAllowed(), nil
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
