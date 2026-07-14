package model

import "time"

type User struct {
	ID        int64
	Username  string
	Nickname  string
	CreatedAt time.Time
	UpdatedAt time.Time
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