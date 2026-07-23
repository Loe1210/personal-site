package xhttp

import (
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/Loe1210/personal-site/internal/xerrors"
)

func Recover() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if value := recover(); value != nil {
				log.Printf("panic recovered: %v", value)
				Fail(c, xerrors.New(xerrors.CodeInternal, "internal error"))
				c.Abort()
			}
		}()
		c.Next(ctx)
	}
}
