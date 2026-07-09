package auth

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	authmodel "github.com/Loe1210/personal-site/biz/model/auth"
	"github.com/Loe1210/personal-site/pkg/response"
)

// Me godoc
// @Summary 获取当前管理员信息
// @Description 通过 JWT 获取当前登录管理员信息
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Body
// @Failure 401 {object} response.Body
// @Router /api/admin/me [get]
func Me(_ context.Context, c *app.RequestContext) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	resp := &authmodel.GetCurrentAdminResponse{
		User: &authmodel.AdminUser{
			ID:       userID.(int64),
			Username: username.(string),
			Nickname: "Administrator",
		},
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}