package auth

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/sessions"

	authmodel "github.com/Loe1210/personal-site/biz/model/auth"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
	dbmodel "github.com/Loe1210/personal-site/biz/dal/db"
)

// Me godoc
// @Summary 获取当前用户信息
// @Description 通过 Session 获取当前登录用户信息
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Body
// @Failure 401 {object} response.Body
// @Router /api/admin/me [get]
func Me(_ context.Context, c *app.RequestContext) {
	session := sessions.Default(c)

	userID := session.Get("user_id")
	username := session.Get("username")
	var user dbmodel.User
	if err := dbmodel.DB.First(&user, userID.(int64)).Error; err != nil {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrorCode, "user not found"))
		return
	}

	if user.Username != username.(string) {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrorCode, "login required"))
		return
	}

	resp := &authmodel.GetCurrentUserResponse{
		User: &authmodel.User{
			ID:        user.ID,
			Username:  user.Username,
			Nickname:  user.Nickname,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}