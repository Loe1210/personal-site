package router

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/Loe1210/personal-site/pkg/response"
)

func registerHealth(h *server.Hertz) {
	h.GET("/health", func(_ context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, response.Success(map[string]any{
			"status": "ok",
		}))
	})
}