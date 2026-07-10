package biz

import (
	"context"

	articlebiz "github.com/Loe1210/personal-site/biz/article"
	authbiz "github.com/Loe1210/personal-site/biz/auth"
	categorybiz "github.com/Loe1210/personal-site/biz/category"
	sitebiz "github.com/Loe1210/personal-site/biz/site"
	tagbiz "github.com/Loe1210/personal-site/biz/tag"
	uploadbiz "github.com/Loe1210/personal-site/biz/upload"
	"github.com/Loe1210/personal-site/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Register(h *server.Hertz) {
	sitebiz.Register(h)
	authbiz.Register(h)
	articlebiz.Register(h)
	categorybiz.Register(h)
	tagbiz.Register(h)
	uploadbiz.Register(h)
	registerHealth(h)
}

func registerHealth(h *server.Hertz) {
	h.GET("/health", func(_ context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, response.Success(map[string]any{
			"status": "ok",
		}))
	})
}
