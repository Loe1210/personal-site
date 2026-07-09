package site

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/Loe1210/personal-site/pkg/response"
)

func Home(_ context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, response.Success(map[string]any{
		"message": "personal site is running",
	}))
}