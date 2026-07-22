package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/Loe1210/personal-site/internal/xerrors"
	"github.com/Loe1210/personal-site/internal/xhttp"
)

type SessionValidator interface {
	ValidateSession(ctx context.Context, sessionID string) error
}

func AuthRequired(validator SessionValidator) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		sessionID := string(c.Cookie("session_id"))
		if sessionID == "" {
			xhttp.Fail(c, xerrors.New(xerrors.CodeAuthLoginRequired, "login required"))
			c.Abort()
			return
		}
		if validator != nil {
			if err := validator.ValidateSession(ctx, sessionID); err != nil {
				xhttp.Fail(c, err)
				c.Abort()
				return
			}
		}
		c.Next(ctx)
	}
}
