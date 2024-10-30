package model

import (
	"gorm.io/gorm"
	"time"
)

// VirtualMachine 虚拟机
type VirtualMachine struct {
	ID              uint           `json:"id" gorm:"not null;unique;primary_key"`
	Name            string         `json:"name" gorm:"comment:'实例名称'"`
	InstanceId      string         `json:"instance_id" gorm:"not null;unique; comment:'实例ID'"`
	UserName        string         `json:"-" gorm:"comment:'用户名'" `
	Password        string         `json:"-" gorm:"comment:'密码'"`
	Port            string         `json:"-" gorm:"comment:'端口'"`
	CPU             int            `json:"cpu" gorm:"comment:'CPU'"`
	Memory          int            `json:"memory" gorm:"comment:'内存'"`
	OS              string         `json:"os" gorm:"comment:'操作系统'"`
	OSType          string         `json:"os_type" gorm:"comment:'系统类型'"`
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

func (v VirtualMachine) TableName() string {
	return "cloud_virtual_machine"
}
