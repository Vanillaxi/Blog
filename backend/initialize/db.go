package initialize

import (
	"MyBlog/global"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() error {
	dsn := global.Config.Database.Dsn

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败：%w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接池失败：%w", err)
	}

	sqlDB.SetMaxIdleConns(global.Config.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(global.Config.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.DB = db

	return nil
}
