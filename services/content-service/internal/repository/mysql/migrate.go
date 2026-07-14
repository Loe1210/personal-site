package mysql

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex"`
	Slug        string `gorm:"type:varchar(100);not null;uniqueIndex"`
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Category{}, &Tag{}, &Article{}, &ArticleTag{})
}
