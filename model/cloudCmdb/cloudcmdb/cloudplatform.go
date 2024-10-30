package cloudcmdb

import (
	model "DYCLOUD/model/cloudCmdb"
	"DYCLOUD/model/common/request"
)

type SearchCloudPlatformParams struct {
	model.CloudPlatform
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type CloudResponse struct {
	CloudPlatform model.CloudPlatform  `json:"cloud_platform"`
	Regions       []model.CloudRegions `json:"regions"`
}

type SearchVirtualMachineParams struct {
	model.VirtualMachine
	PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type SearchLoadBalancerParams struct {
	model.LoadBalancer
	PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type SearchRDSParams struct {
	model.RDS
	PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type PageInfo struct {
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页大小
	Keyword  string `json:"keyword" form:"keyword"`   //关键字
	Field    string `json:"field" form:"field"`       // 搜索字段
}
