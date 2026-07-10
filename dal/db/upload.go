package db

import "time"

type UploadFile struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	FileName  string    `gorm:"type:varchar(255);not null"`
	FileURL   string    `gorm:"type:varchar(255);not null"`
	FilePath  string    `gorm:"type:varchar(512);not null"`
	MimeType  string    `gorm:"type:varchar(128);not null"`
	Size      int64     `gorm:"not null"`
	BizType   string    `gorm:"type:varchar(64);not null;default:'common'"`
	CreatedAt time.Time
}