package models

import (
	"MyBlog/global"
	"time"
)

type Tag struct {
	ID           uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	TagName      string    `gorm:"size:255;not null" json:"tag_name"`
	ArticleCount uint32    `gorm:"column:article_count" json:"article_count"`
	Status       int8      `gorm:"column:status" json:"status"`
	Slug         string    `gorm:"size:255;not null" json:"slug"`
	Sort         uint32    `gorm:"column:sort" json:"sort"`
	CreateTime   time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"update_time"`
}

func (Tag) TableName() string {
	return "tag"
}

func CheckTagsExist(tagIDs []uint32) (bool, error) {
	if len(tagIDs) == 0 {
		return true, nil
	}

	var count int64
	err := global.DB.Model(&Tag{}).
		Where("id in (?)", tagIDs).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == int64(len(tagIDs)), nil
}

func CreateTag(tag *Tag) error {
	return global.DB.Create(tag).Error
}

// 前台展示：只查启用的标签
func GetTags() ([]Tag, error) {
	var tag []Tag
	err := global.DB.
		Where("status = ?", 1).
		Order("sort desc,id desc").
		Find(&tag).Error
	return tag, err
}

// 后台管理：查全部标签
func GetAdminTags() ([]Tag, error) {
	var tag []Tag
	err := global.DB.
		Order("sort desc,id desc").
		Find(&tag).Error
	return tag, err
}

func UpdateTagStatus(id uint32, status int8) error {
	return global.DB.Model(&Tag{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// 排序
func UpdateTagSort(id uint32, sort uint32) error {
	return global.DB.Model(&Tag{}).
		Where("id=?", id).
		Update("sort", sort).Error
}
