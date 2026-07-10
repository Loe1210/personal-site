package auth

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/sessions"

	authmodel "github.com/Loe1210/personal-site/biz/model/auth"
	authservice "github.com/Loe1210/personal-site/biz/service/auth"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
)

// Login godoc
// @Summary 用户登录
// @Description 使用用户账号密码登录并获取登录态
// @Tags auth
// @Accept json
// @Produce json
// @Param body body auth.UserLoginRequest true "登录请求"
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Failure 401 {object} response.Body
// @Router /api/admin/login [post]
func Login(ctx context.Context, c *app.RequestContext) {
	var req authmodel.UserLoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	resp, err := authservice.Login(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", resp.User.ID)
	session.Set("username", resp.User.Username)

	if err := session.Save(); err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.ErrorCode, "save session failed"))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}