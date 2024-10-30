package system

import (
	"DYCLOUD/global"
)

type SysApi struct {
	global.DYCLOUD_MODEL
	Path        string `json:"path" gorm:"comment:api路径"`             // api路径
	Description string `json:"description" gorm:"comment:api中文描述"`    // api中文描述
	ApiGroup    string `json:"apiGroup" gorm:"comment:api组"`          // api组
	Method      string `json:"method" gorm:"default:POST;comment:方法"` // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
}

func (SysApi) TableName() string {
	return "sys_apis"
}

type SysIgnoreApi struct {
	global.DYCLOUD_MODEL
	Path   string `json:"path" gorm:"comment:api路径"`             // api路径
	Method string `json:"method" gorm:"default:POST;comment:方法"` // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
	Flag   bool   `json:"flag" gorm:"-"`                         // 是否忽略
}

func (SysIgnoreApi) TableName() string {
	return "sys_ignore_apis"
}
