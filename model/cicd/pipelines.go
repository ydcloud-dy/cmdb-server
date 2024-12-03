package cicd

import "DYCLOUD/global"

type Pipelines struct {
	global.DYCLOUD_MODEL
	Name           string  `json:"name" form:"name"`
	AppName        string  `json:"app_name" form:"app_name"`
	EnvName        string  `json:"env_name" form:"env_name"`
	BuildScript    string  `json:"build_script" form:"build_script"`
	K8SNamespace   string  `json:"k8s_namespace" form:"k8s_namespace" gorm:"column:k8s_namespace"`
	K8SClusterName string  `json:"k8s_cluster_name" form:"k8s_cluster_name" gorm:"column:k8s_cluster_name"`
	BaseImage      string  `json:"base_image" form:"base_image"`
	DockerfilePath string  `json:"dockerfile_path" form:"dockerfile_path"`
	ImageName      string  `json:"image_name" form:"image_name"`
	ImageTag       string  `json:"image_tag" form:"image_tag"`
	RegistryURL    string  `json:"registry_url" form:"registry_url"`
	RegistryUser   string  `json:"registry_user" form:"registry_user"`
	RegistryPass   string  `json:"registry_pass" form:"registry_pass"`
	GitUrl         string  `json:"git_url" form:"git_url"`
	GitBranch      string  `json:"git_branch" form:"git_branch"`
	GitCommitId    string  `json:"git_commit_id" form:"git_commit_id"`
	Stages         []Stage `json:"stages" gorm:"foreignKey:PipelineID"`
	CreatedBy      uint    `gorm:"column:created_by;comment:创建者"`
	CreatedName    string  `gorm:"column:created_name;comment:创建者名字"`
	UpdatedName    string  `gorm:"column:updated_name;comment:修改者名字"`
	UpdatedBy      uint    `gorm:"column:updated_by;comment:更新者"`
	DeletedBy      uint    `gorm:"column:deleted_by;comment:删除者"`
}
type Stage struct {
	ID         uint   `gorm:"primaryKey"`
	PipelineID uint   `json:"pipeline_id"`
	Name       string `json:"name"`
	TaskList   []Task `json:"task_list" gorm:"foreignKey:StageID"`
}

func (s *Stage) TableName() string {
	return "cicd_pipelines_stages"
}

func (p *Pipelines) TableName() string {
	return "cicd_pipelines"
}

type Task struct {
	ID           uint   `gorm:"primaryKey"`
	StageID      uint   `json:"stage_id"`
	Name         string `json:"name"`
	Branch       string `json:"branch"`
	Image        string `json:"image"`
	Plugin       string `json:"plugin"` // 任务类型字段
	Script       string `json:"script"`
	SpatialName  string `json:"spatial_name"`
	Warehouse    string `json:"warehouse"`
	MirrorTag    string `json:"mirror_tag"`
	Dockerfile   string `json:"dockerfile"`
	ContextPath  string `json:"context_path"`
	ProductName  string `json:"product_name"`
	ProductPath  string `json:"product_path"`
	Version      string `json:"version"`
	Resource     string `json:"resource"`
	YamlResource string `json:"yaml_resource"`
	GoalResource string `json:"goal_resource"`
	OpenScript   string `json:"open_script"`
	PipelineID   uint   `json:"pipeline_id"`
}

func (t *Task) TableName() string {
	return "cicd_pipelines_tasks"
}
