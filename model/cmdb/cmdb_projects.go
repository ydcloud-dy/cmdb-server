// 自动生成模板CmdbProjects
package cmdb

import (
	"DYCLOUD/global"
)

// cmdbProjects表 结构体  CmdbProjects
type CmdbProjects struct {
	global.DYCLOUD_MODEL
	Name      string `json:"name" form:"name" gorm:"column:name;comment:项目名称;size:255;" binding:"required"`               //项目名称
	Principal string `json:"principal" form:"principal" gorm:"column:principal;comment:负责人;size:255;" binding:"required"` //负责人
	Note      string `json:"note" form:"note" gorm:"column:note;comment:备注;size:255;"`                                    //备注
	CreatedBy uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName cmdbProjects表 CmdbProjects自定义表名 cmdb_projects
func (CmdbProjects) TableName() string {
	return "cmdb_projects"
}
