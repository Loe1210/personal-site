package auth

import (
	"context"
	"errors"

	dbmodel "github.com/Loe1210/personal-site/biz/dal/db"
	authmodel "github.com/Loe1210/personal-site/biz/model/auth"
	"golang.org/x/crypto/bcrypt"
)

func Login(_ context.Context, req *authmodel.UserLoginRequest) (*authmodel.UserLoginResponse, error) {
	var user dbmodel.User

	if err := dbmodel.DB.Where("username = ? AND status = ?", req.Username, "active").First(&user).Error; err != nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return &authmodel.UserLoginResponse{
		User: &authmodel.User{
			ID:        user.ID,
			Username:  user.Username,
			Nickname:  user.Nickname,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}