package cicd

import "DYCLOUD/global"

// Environment
// @Description: 应用发布环境
type Environment struct {
	global.DYCLOUD_MODEL
	Name      string `json:"name" form:"name"`
	Key       string `json:"key" form:"key"`
	Desc      string `json:"desc" form:"desc"`
	CreatedBy uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint   `gorm:"column:deleted_by;comment:删除者"`
}

func (e *Environment) TableName() string {
	return "cicd_environments"
}
