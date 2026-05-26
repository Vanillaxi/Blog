package models

import (
	"MyBlog/global"
	"time"
)

type Category struct {
	ID           uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryName string    `gorm:"size:255;not null" json:"category_name"`
	Sort         uint32    `gorm:"not null;default:0" json:"sort"`
	Slug         string    `gorm:"size:255;not null" json:"slug"`
	Status       int       `gorm:"not null;default:1" json:"status"`
	ArticleCount uint32    `gorm:"not null;default:0" json:"article_count"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

func (Category) TableName() string {
	return "category"
}

func CheckCategoryExist(categoryID uint32) (bool, error) {
	var count int64
	err := global.DB.Model(&Category{}).
		Where("id = ? and status=?", categoryID, 1).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// 前台展示：只查启用的分类
func GetCategories() ([]Category, error) {
	var categories []Category
	err := global.DB.
		Where("status = ?", 1).
		Order("sort desc,id desc").
		Find(&categories).Error
	return categories, err
}

// 后台管理：查全部分类
func GetAdminCategories() ([]Category, error) {
	var categories []Category
	err := global.DB.
		Order("sort desc,id desc").
		Find(&categories).Error
	return categories, err
}

func CreateCategory(category *Category) error {
	return global.DB.Create(category).Error
}

func UpdateCategoryStatus(id uint32, status int8) error {
	return global.DB.Model(&Category{}).
		Where("id=?", id).
		Update("status", status).Error
}

func UpdateCategorySort(id uint32, sort uint32) error {
	return global.DB.Model(&Category{}).
		Where("id=?", id).
		Update("sort", sort).Error
}
