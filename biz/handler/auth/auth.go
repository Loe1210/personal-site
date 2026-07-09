package auth

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	authmodel "github.com/Loe1210/personal-site/biz/model/auth"
	authservice "github.com/Loe1210/personal-site/biz/service/auth"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
)

// Login godoc
// @Summary 管理员登录
// @Description 使用管理员账号密码登录并获取 JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param body body auth.AdminLoginRequest true "登录请求"
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Failure 401 {object} response.Body
// @Router /api/admin/login [post]
func Login(ctx context.Context, c *app.RequestContext) {
	var req authmodel.AdminLoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	resp, err := authservice.Login(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}