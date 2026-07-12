package db

import "time"

type Article struct {
	ID          int64      `gorm:"primaryKey;autoIncrement"`
	Title       string     `gorm:"type:varchar(255);not null"`
	Slug        string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	Summary     string     `gorm:"type:text"`
	ContentMd   string     `gorm:"type:longtext"`
	ContentHTML string     `gorm:"type:longtext"`
	CoverImage  string     `gorm:"type:varchar(255)"`
	CategoryID  int64      `gorm:"default:0"`
	Status      string     `gorm:"type:varchar(32);index;default:'draft'"`
	IsTop       int        `gorm:"type:tinyint(1);default:0;index"`
	AuthorID    int64      `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt *time.Time `gorm:"default:null"`
}