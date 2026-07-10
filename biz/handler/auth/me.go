package auth

import (
	"context"

	dbmodel "github.com/Loe1210/personal-site/biz/dal/db"
	authmodel "github.com/Loe1210/personal-site/biz/model/auth"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/sessions"
)

func Me(_ context.Context, c *app.RequestContext) {
	session := sessions.Default(c)

	userIDValue := session.Get("user_id")
	usernameValue := session.Get("username")
	if userIDValue == nil || usernameValue == nil {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.Unauthorized.Code, "login required"))
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.Unauthorized.Code, "invalid session"))
		return
	}
	username, ok := usernameValue.(string)
	if !ok {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.Unauthorized.Code, "invalid session"))
		return
	}

	var user dbmodel.User
	if err := dbmodel.DB.First(&user, userID).Error; err != nil {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrorCode, "user not found"))
		return
	}
	if user.Username != username {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.Unauthorized.Code, "login required"))
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
