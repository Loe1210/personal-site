package auth

import "context"

type Context struct {
	UserID   int64
	Username string
	Roles    []string
}

type Client interface {
	ValidateSession(ctx context.Context, sessionID string) (*Context, error)
	CheckPermission(ctx context.Context, userID int64, code string) (bool, error)
}
