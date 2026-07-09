package db

import "time"

type Tag struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	Slug        string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ArticleTag struct {
	ArticleID int64 `gorm:"primaryKey"`
	TagID     int64 `gorm:"primaryKey"`
}