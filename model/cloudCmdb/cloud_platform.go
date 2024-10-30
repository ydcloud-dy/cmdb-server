package model

import (
	"gorm.io/gorm"
	"time"
)

// CloudPlatform 云厂商
type CloudPlatform struct {
	ID              uint           `json:"id" gorm:"not null;unique;primary_key"`
	Name            string         `json:"name" form:"name" gorm:"comment:厂商名称"`
	AccessKeyId     string         `json:"access_key_id"`
	AccessKeySecret string         `json:"access_key_secret"`
	Platform        string         `json:"platform"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c CloudPlatform) TableName() string {
	return "cloud_platform"
}
