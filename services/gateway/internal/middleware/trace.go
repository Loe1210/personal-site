package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func Trace() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)
	}
}
