package application

import (
	"context"

	"github.com/Loe1210/personal-site/pkg/xauth"
	"github.com/Loe1210/personal-site/services/auth-service/internal/domain"
)

type User = domain.User

type UserRepository interface {
	Login(ctx context.Context, username, password string) (*User, []string, error)
	GetByID(ctx context.Context, userID int64) (*User, error)
	HasPermission(ctx context.Context, userID int64, code string) (bool, error)
}

type SessionBundle struct {
	SessionID  string
	CookieName string
	ExpiresAt  string
	Username   string
}

type AuthContext struct {
	UserID   int64
	Username string
	Roles    []string
}

type Service struct {
	users UserRepository
}

func NewAuthService(users UserRepository) *Service {
	return &Service{users: users}
}

func (s *Service) CreateSession(ctx context.Context, username string, password string) (*SessionBundle, error) {
	user, roles, err := s.users.Login(ctx, username, password)
	if err != nil {
		return nil, err
	}

	bundle, err := xauth.CreateSessionBundle(user.ID, user.Username, roles)
	if err != nil {
		return nil, err
	}

	return &SessionBundle{
		SessionID:  bundle.SessionID,
		CookieName: bundle.CookieName,
		ExpiresAt:  bundle.ExpiresAt.Format(timeLayout),
		Username:   user.Username,
	}, nil
}

func (s *Service) ValidateSession(_ context.Context, sessionID string) (*AuthContext, error) {
	claims, err := xauth.ParseSession(sessionID)
	if err != nil {
		return nil, err
	}

	return &AuthContext{
		UserID:   claims.UserID,
		Username: claims.Username,
		Roles:    append([]string(nil), claims.Roles...),
	}, nil
}

func (s *Service) CheckPermission(ctx context.Context, userID int64, code string) (bool, error) {
	return s.users.HasPermission(ctx, userID, code)
}

func (s *Service) GetCurrentUser(ctx context.Context, sessionID string) (*User, error) {
	authContext, err := s.ValidateSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return s.users.GetByID(ctx, authContext.UserID)
}

const timeLayout = "2006-01-02 15:04:05"
