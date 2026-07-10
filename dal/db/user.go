package db

import "time"

type User struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	Username     string `gorm:"type:varchar(64);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Nickname     string `gorm:"type:varchar(64);not null"`
	Type         string `gorm:"type:varchar(32);index;default:'admin'"`
	Status       string `gorm:"type:varchar(32);index;default:'active'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}