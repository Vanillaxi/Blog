package global

import (
	"MyBlog/config"

	"gorm.io/gorm"
)

var (
	Config *config.Config
	DB     *gorm.DB
)
