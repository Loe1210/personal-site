package mysql

import (
	"fmt"

	"github.com/Loe1210/personal-site/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open(cfg configs.MySQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
