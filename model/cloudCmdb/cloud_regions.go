package model

import (
	"gorm.io/gorm"
	"time"
)

type Regions struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	RegionId   string `json:"region_id"`
	RegionName string `json:"region_name"`
}

// CloudRegions 云资产地域信息
type CloudRegions struct {
	ID              uint           `json:"id" gorm:"not null;unique;primary_key"`
	Name            string         `json:"name"`
	RegionId        string         `json:"region_id" gorm:"not null;unique; comment:'区域ID'"`
	RegionName      string         `json:"region_name" gorm:"comment:'区域名称'"`
	CloudPlatformId uint           `json:"cloud_platform_id"`
	CloudPlatform   CloudPlatform  `json:"cloud_platform" gorm:"ForeignKey:CloudPlatformId"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r Regions) TableName() string {
	return "cloud_regions"
}
