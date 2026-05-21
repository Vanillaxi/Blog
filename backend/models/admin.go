package models

import (
	"MyBlog/global"
	"time"
)

type Admin struct {
	ID         uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username   string    `gorm:"not null;unique" json:"username"`
	Password   string    `gorm:"not null" json:"password"` //	密码通常在JSON中隐藏
	Nickname   string    `gorm:"not null" json:"nickname"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
}

func GetAdminByUsername(username string) (Admin, error) {
	var admin Admin
	err := global.DB.Where("username = ?", username).First(&admin).Error
	return admin, err
}

func GetAdminByID(id uint64) (Admin, error) {
	var admin Admin
	err := global.DB.Where("id = ?", id).First(&admin).Error
	return admin, err
}

func CreateAdmin(m *Admin) error {
	return global.DB.Create(m).Error
}

func UpdateAdminName(id uint32, newUsername string) error {
	return global.DB.Model(&Admin{}).Where("id = ?", id).Update("username", newUsername).Error
}

func UpdatePassword(id uint32, newPassword string) error {
	return global.DB.Model(&Admin{}).Where("id = ?", id).Update("password", newPassword).Error
}
