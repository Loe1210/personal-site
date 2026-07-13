package session

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
	"github.com/Loe1210/personal-site/pkg/xauth"
)

func AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		sessionID := xauth.SessionIDFromRequest(c)
		claims, err := xauth.ParseSession(sessionID)
		if err != nil {
			response.WriteErrorMessage(c, errno.Unauthorized, "login required")
			c.Abort()
			return
		}

		xauth.SetClaims(c, claims)
		c.Next(ctx)
	}
}
