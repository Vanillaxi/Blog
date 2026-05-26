package models

import "time"

type VisitorIdentity struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Email      string    `gorm:"size:100;not null;uniqueIndex:uk_visitor_identity_email" json:"email"`
	Nickname   string    `gorm:"size:50;not null" json:"nickname"`
	Avatar     string    `gorm:"size:255" json:"avatar"`
	FirstIP    string    `gorm:"column:first_ip;size:64" json:"first_ip"`
	LastIP     string    `gorm:"column:last_ip;size:64" json:"last_ip"`
	UserAgent  string    `gorm:"size:500" json:"user_agent"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

func (VisitorIdentity) TableName() string {
	return "visitor_identity"
}
