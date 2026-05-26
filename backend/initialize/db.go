package initialize

import (
	"MyBlog/global"
	"MyBlog/models"
	"fmt"
	"strings"
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
	if err := db.AutoMigrate(&models.VisitorIdentity{}); err != nil {
		return fmt.Errorf("初始化访客身份表失败：%w", err)
	}
	if !db.Migrator().HasColumn(&models.FriendLink{}, "Description") {
		if err := db.Migrator().AddColumn(&models.FriendLink{}, "Description"); err != nil {
			return fmt.Errorf("初始化友链描述字段失败：%w", err)
		}
	}
	if !db.Migrator().HasColumn(&models.Comment{}, "IPLocation") {
		if err := db.Migrator().AddColumn(&models.Comment{}, "IPLocation"); err != nil {
			return fmt.Errorf("初始化评论地区字段失败：%w", err)
		}
	}
	if err := ensureArticleContentColumn(db); err != nil {
		return err
	}

	return nil
}

func ensureArticleContentColumn(db *gorm.DB) error {
	columnTypes, err := db.Migrator().ColumnTypes(&models.Article{})
	if err != nil {
		return fmt.Errorf("检查文章正文字段失败：%w", err)
	}

	for _, columnType := range columnTypes {
		if columnType.Name() != "content" {
			continue
		}

		databaseType := columnType.DatabaseTypeName()
		if strings.EqualFold(databaseType, "LONGTEXT") || strings.EqualFold(databaseType, "TEXT") {
			return nil
		}

		if err := db.Migrator().AlterColumn(&models.Article{}, "Content"); err != nil {
			return fmt.Errorf("修正文章正文字段类型失败：%w", err)
		}
		return nil
	}

	return nil
}
