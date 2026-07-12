package service

import (
	"context"

	dbmodel "github.com/Loe1210/personal-site/dal/db"
	authmodel "github.com/Loe1210/personal-site/biz/model/auth"
	"github.com/Loe1210/personal-site/pkg/errno"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(_ context.Context, req *authmodel.UserLoginRequest) (*authmodel.UserLoginResponse, error) {
	var user dbmodel.User

	if err := dbmodel.DB.Where("username = ? AND status = ?", req.Username, "active").First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.InvalidCredentials
		}
		return nil, errno.Internal
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errno.InvalidCredentials
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
