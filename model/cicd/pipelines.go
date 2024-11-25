package cicd

import "DYCLOUD/global"

type Pipelines struct {
	global.DYCLOUD_MODEL
	Name           string `json:"name" form:"name"`
	AppName        string `json:"app_name" form:"app_name"`
	EnvName        string `json:"env_name" form:"env_name"`
	BuildScript    string `json:"build_script" form:"build_script"`
	K8SNamespace   string `json:"k8s_namespace" form:"k8s_namespace" gorm:"column:k8s_namespace"`
	K8SClusterName string `json:"k8s_cluster_name" form:"k8s_cluster_name" gorm:"k8s_cluster_name"`
	// 展开的 Docker 和 Registry 字段
	BaseImage      string `json:"base_image" form:"base_image"`
	DockerfilePath string `json:"dockerfile_path" form:"dockerfile_path"`
	ImageName      string `json:"image_name" form:"image_name"`
	ImageTag       string `json:"image_tag" form:"image_tag"`
	RegistryURL    string `json:"registry_url" form:"registry_url"`
	RegistryUser   string `json:"registry_user" form:"registry_user"`
	RegistryPass   string `json:"registry_pass" form:"registry_pass"`

	// Repository 字段
	GitUrl      string `json:"git_url" form:"git_url"`
	GitBranch   string `json:"git_branch" form:"git_branch"`
	GitCommitId string `json:"git_commit_id" form:"git_commit_id"`
}

func (p *Pipelines) TableName() string {
	return "cicd_pipelines"
}
