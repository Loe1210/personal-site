package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Loe1210/personal-site/pkg/configs"
)

var DB *gorm.DB

func Init() error {
	cfg := configs.LoadMySQLConfig()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := database.AutoMigrate(&Article{}, &Category{}, &Tag{}, &ArticleTag{}); err != nil {
		return err
	}

	DB = database

	return nil
}