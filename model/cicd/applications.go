package cicd

import (
	"DYCLOUD/global"
	"encoding/json"
)

type Applications struct {
	global.DYCLOUD_MODEL
	Name         string `json:"name" form:"name"`
	FullName     string `json:"full_name" form:"full_name"`
	Language     string `json:"language" form:"language"`
	BuildPath    string `json:"build_path" form:"build_path"`
	Dockerfile   string `json:"dockerfile" form:"dockerfile"`
	RepoId       int    `json:"repo_id" form:"repo_id"`
	Path         string `json:"path" form:"path"`
	CompileEnvId int    `json:"compile_env_id" form:"compile_env"`
	CreatedBy    uint   `gorm:"column:created_by;comment:创建者"`
	CreatedName  string `gorm:"column:created_name;comment:创建者名字"`
	UpdatedName  string `gorm:"column:updated_name;comment:修改者名字"`
	UpdatedBy    uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy    uint   `gorm:"column:deleted_by;comment:删除者"`
}

func (a *Applications) TableName() string {
	return "cicd_applications"
}
func (a *Applications) String() string {
	data, err := json.Marshal(&a)
	if err != nil {
		return ""
	}
	return string(data)
}

type AppBranch struct {
	global.DYCLOUD_MODEL
	AppID       int    `orm:"column(app_id);" json:"app_id"`
	BranchName  string `orm:"column(branch_name);size(64)" json:"branch_name"`
	Path        string `orm:"column(path);size(256)" json:"path"`
	CreatedBy   uint   `gorm:"column:created_by;comment:创建者"`
	CreatedName string `gorm:"column:created_name;comment:创建者名字"`
	UpdatedName string `gorm:"column:updated_name;comment:修改者名字"`
	UpdatedBy   uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy   uint   `gorm:"column:deleted_by;comment:删除者"`
}

func (a *AppBranch) TableName() string {
	return "cicd_app_branchs"
}
