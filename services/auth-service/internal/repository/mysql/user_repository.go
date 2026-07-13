package mysql

import (
	"context"
	"errors"

	"github.com/Loe1210/personal-site/services/auth-service/internal/application"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRecord struct {
	ID           int64
	Username     string
	PasswordHash string
	Nickname     string
	Status       string
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Login(ctx context.Context, username, password string) (*application.User, []string, error) {
	var record userRecord
	if err := r.db.WithContext(ctx).Table("users").Where("username = ? AND status = ?", username, "active").First(&record).Error; err != nil {
		return nil, nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(record.PasswordHash), []byte(password)); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}
	roles, err := r.rolesByUserID(ctx, record.ID)
	if err != nil {
		return nil, nil, err
	}
	return &application.User{ID: record.ID, Username: record.Username, Nickname: record.Nickname}, roles, nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID int64) (*application.User, error) {
	var record userRecord
	if err := r.db.WithContext(ctx).Table("users").First(&record, userID).Error; err != nil {
		return nil, err
	}
	return &application.User{ID: record.ID, Username: record.Username, Nickname: record.Nickname}, nil
}

func (r *UserRepository) HasPermission(ctx context.Context, userID int64, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("user_roles").
		Joins("JOIN role_permissions ON user_roles.role_id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("user_roles.user_id = ? AND permissions.code = ?", userID, code).
		Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) rolesByUserID(ctx context.Context, userID int64) ([]string, error) {
	var roles []string
	err := r.db.WithContext(ctx).Table("roles").
		Select("roles.code").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Order("roles.id ASC").
		Scan(&roles).Error
	return roles, err
}
