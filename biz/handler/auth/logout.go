package auth

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/sessions"

	authmodel "github.com/Loe1210/personal-site/biz/model/auth"
	"github.com/Loe1210/personal-site/pkg/response"
)

// Logout godoc
// @Summary 用户登出
// @Description 清除当前登录 Session
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Body
// @Router /api/admin/logout [post]
func Logout(_ context.Context, c *app.RequestContext) {
	session := sessions.Default(c)
	session.Clear()

	if err := session.Save(); err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(10000, "logout failed"))
		return
	}

	c.JSON(consts.StatusOK, response.Success(&authmodel.LogoutResponse{
		Message: "logout success",
	}))
}