package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type SessionValidator interface {
	ValidateSession(ctx context.Context, sessionID string) error
}

func AuthRequired(validator SessionValidator) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		sessionID := string(c.Cookie("session_id"))
		if sessionID == "" {
			c.JSON(consts.StatusUnauthorized, map[string]any{"code": 10002, "message": "login required"})
			c.Abort()
			return
		}
		if validator != nil {
			if err := validator.ValidateSession(ctx, sessionID); err != nil {
				c.JSON(consts.StatusUnauthorized, map[string]any{"code": 10002, "message": "login expired"})
				c.Abort()
				return
			}
		}
		c.Next(ctx)
	}
}
