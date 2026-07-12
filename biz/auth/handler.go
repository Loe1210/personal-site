package auth

import (
	"context"

	authmodel "github.com/Loe1210/personal-site/biz/model/auth"
	dbmodel "github.com/Loe1210/personal-site/dal/db"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
	authservice "github.com/Loe1210/personal-site/service"
	"github.com/cloudwego/hertz/pkg/app"
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
		if appErr, ok := err.(*errno.AppError); ok {
			response.WriteError(c, appErr)
			return
		}
		response.WriteError(c, errno.Internal)
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", resp.User.ID)
	session.Set("username", resp.User.Username)

	if err := session.Save(); err != nil {
		response.WriteError(c, errno.Internal)
		return
	}

	response.WriteSuccess(c, resp)
}

func Me(_ context.Context, c *app.RequestContext) {
	session := sessions.Default(c)

	userIDValue := session.Get("user_id")
	usernameValue := session.Get("username")
	if userIDValue == nil || usernameValue == nil {
		response.WriteError(c, errno.Unauthorized)
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		response.WriteError(c, errno.Unauthorized)
		return
	}
	username, ok := usernameValue.(string)
	if !ok {
		response.WriteError(c, errno.Unauthorized)
		return
	}

	var user dbmodel.User
	if err := dbmodel.DB.First(&user, userID).Error; err != nil {
		response.WriteError(c, errno.Unauthorized)
		return
	}
	if user.Username != username {
		response.WriteError(c, errno.Unauthorized)
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

	response.WriteSuccess(c, resp)
}

func Logout(_ context.Context, c *app.RequestContext) {
	session := sessions.Default(c)
	session.Clear()

	if err := session.Save(); err != nil {
		response.WriteError(c, errno.Internal)
		return
	}

	response.WriteSuccess(c, &authmodel.LogoutResponse{
		Message: "logout success",
	})
}
