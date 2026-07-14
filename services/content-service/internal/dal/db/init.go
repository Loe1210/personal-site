package db

import (
	"fmt"
	"time"

	"github.com/Loe1210/personal-site/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type categoryMigrationModel struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(100);not null;uniqueIndex:uk_categories_name"`
	Slug        string    `gorm:"type:varchar(100);not null;uniqueIndex:uk_categories_slug"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
}

func (categoryMigrationModel) TableName() string {
	return "categories"
}

func Open(cfg configs.MySQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&categoryMigrationModel{}, &Tag{}, &Article{}, &ArticleTag{})
}
