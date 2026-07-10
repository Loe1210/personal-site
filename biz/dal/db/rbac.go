package db

import "time"

type Role struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(64);not null"`
	Code        string `gorm:"type:varchar(64);uniqueIndex;not null"`
	Description string `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Permission struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(64);not null"`
	Code        string `gorm:"type:varchar(64);uniqueIndex;not null"`
	Resource    string `gorm:"type:varchar(64);index;not null"`
	Action      string `gorm:"type:varchar(64);index;not null"`
	Description string `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserRole struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64 `gorm:"uniqueIndex:idx_user_role;not null"`
	RoleID    int64 `gorm:"uniqueIndex:idx_user_role;not null"`
	CreatedAt time.Time
}

type RolePermission struct {
	ID           int64 `gorm:"primaryKey;autoIncrement"`
	RoleID       int64 `gorm:"uniqueIndex:idx_role_permission;not null"`
	PermissionID int64 `gorm:"uniqueIndex:idx_role_permission;not null"`
	CreatedAt    time.Time
}