package model

import (
	"gorm.io/gorm"
	"time"
)

// RDS 数据库
type RDS struct {
	ID              uint           `json:"id" gorm:"not null;unique;primary_key"`
	Name            string         `json:"name" gorm:"comment:'实例名称'"`
	InstanceId      string         `json:"instance_id" gorm:"not null;unique; comment:'实例ID'"`
	PrivateAddr     string         `json:"private_addr" gorm:"comment:'私网地址'"`
	PublicAddr      string         `json:"public_addr" gorm:"comment:'公网地址'"`
	Region          string         `json:"region" gorm:"comment:'区域ID'" `
	RegionName      string         `json:"region_name" gorm:"comment:'区域名称'" `
	Status          string         `json:"status" gorm:"comment:'状态'"`
	CreationTime    string         `json:"creation_time" gorm:"comment:'创建时间'"`
	ExpiredTime     string         `json:"expired_time" gorm:"comment:'到期时间'"`
	CloudPlatformId uint           `json:"cloud_platform_id"`
	CloudPlatform   CloudPlatform  `json:"cloud_platform" gorm:"ForeignKey:CloudPlatformId"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

func (l RDS) TableName() string {
	return "cloud_rds"
}
