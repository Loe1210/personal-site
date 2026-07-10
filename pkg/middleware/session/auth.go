package session

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sessions"

	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
)

func AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		session := sessions.Default(c)

		userID := session.Get("user_id")
		username := session.Get("username")

		if userID == nil || username == nil {
			response.WriteErrorMessage(c, errno.Unauthorized, "login required")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("username", username)

		c.Next(ctx)
	}
}
