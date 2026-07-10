package auth

import (
	"context"

	dbmodel "github.com/Loe1210/personal-site/dal/db"
	authmodel "github.com/Loe1210/personal-site/biz/model/auth"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
	authservice "github.com/Loe1210/personal-site/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/sessions"
)

func Login(ctx context.Context, c *app.RequestContext) {
	req, err := bindLoginRequest(c)
	if err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := authservice.Login(ctx, req)
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
