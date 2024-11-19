package cicd

import (
	"DYCLOUD/global"
	"encoding/json"
	"gorm.io/gorm"
	"time"
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

type Option struct {
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Username string `json:"username"`
}

type Developer struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"ID"`
	CreatedAt   time.Time      // 创建时间
	UpdatedAt   time.Time      // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
	AppID       uint           `gorm:"index" json:"app_id"`
	RoleType    string         `json:"role_type"` // develop 或 owner
	Label       string         `json:"label"`
	Value       string         `json:"value"`
	Key         string         `json:"key"`
	Option      *Option        `gorm:"-" json:"option"` // 忽略 GORM 的处理，只用于解析前端数据
	Avatar      string         `json:"avatar"`
	Nickname    string         `json:"nickname"`
	Username    string         `json:"username"`
	OriginLabel string         `json:"origin_label"`
}

func (d *Developer) String() string {
	data, err := json.Marshal(&d)
	if err != nil {
		return ""
	}
	return string(data)
}
func (d *Developer) TableName() string {
	return "cicd_app_developers"
}

type Env struct {
	global.DYCLOUD_MODEL
	ClusterName string `json:"clusterName"`
	ClusterID   int    `json:"clusterId"`
	EnvCode     string `json:"envCode"`
	EnvName     string `json:"envName"`
	Namespace   string `json:"namespace"`
	AppID       uint   `gorm:"index" json:"app_id"`
}

func (e *Env) TableName() string {
	return "cicd_app_envs"
}

type App struct {
	global.DYCLOUD_MODEL
	GitRepo     string      `json:"gitRepo"`
	Branch      string      `json:"branch"`
	AppName     string      `json:"appName"`
	AppCode     string      `json:"appCode"`
	AppDesc     string      `json:"appDesc"`
	Develop     []Developer `gorm:"-" json:"develop"` // 禁用自动级联插入
	Owner       []Developer `gorm:"-" json:"owner"`   // 禁用自动级联插入
	Envs        []Env       `gorm:"-" json:"envs"`    // 禁用自动级联插入
	Language    string      `json:"language"`
	CreatedBy   uint        `gorm:"column:created_by;comment:创建者"`
	CreatedName string      `gorm:"column:created_name;comment:创建者名字"`
	UpdatedName string      `gorm:"column:updated_name;comment:修改者名字"`
	UpdatedBy   uint        `gorm:"column:updated_by;comment:更新者"`
	DeletedBy   uint        `gorm:"column:deleted_by;comment:删除者"`
}

func (a *App) TableName() string {
	return "cicd_apps"
}

type AppRequestBody struct {
	Envs []Env `json:"envs"`
	App  App   `json:"app"`
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
