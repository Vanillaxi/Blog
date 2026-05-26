package models

import (
	"MyBlog/global"
	"time"
)

type FriendLink struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	URL         string    `gorm:"size:255;not null" json:"url"`
	Logo        string    `gorm:"size:255" json:"logo"`
	Description string    `gorm:"size:255" json:"description"`
	Sort        int       `gorm:"not null;default:0" json:"sort"`
	Status      int       `gorm:"not null;default:1" json:"status"` //1正常，0隐藏
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

func (FriendLink) TableName() string {
	return "friend_link"
}

// 查询前台显示的友链
func GetActiveFriendLinks(pageNum, pageSize int) ([]FriendLink, int64, error) {
	var links []FriendLink
	var total int64

	db := global.DB.Model(&FriendLink{}).Where("status=?", 1)
	db.Count(&total)

	err := db.
		Order("sort DESC,id DESC").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&links).Error

	return links, total, err
}

// 后台分页查询
func GetAdminFriendLinks(pageNum, pageSize int) ([]FriendLink, int64, error) {
	var links []FriendLink
	var total int64

	db := global.DB.Model(&FriendLink{})
	db.Count(&total)

	err := db.Order("sort DESC,id DESC").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&links).Error
	return links, total, err

}

// 创建友链
func CreateFriendLink(link *FriendLink) error {
	return global.DB.Create(link).Error
}

// 更新友链
func UpdateFriendLink(link *FriendLink) error {
	return global.DB.Save(link).Error
}

// 删除友链
func DeleteFriendLink(id uint64) error {
	return global.DB.Delete(&FriendLink{ID: id}).Error
}
