package cicd

import "DYCLOUD/global"

type Pipelines struct {
	global.DYCLOUD_MODEL
	Name         string     `json:"name" form:"name"`
	AppName      string     `json:"app_name" form:"app_name"`
	EnvName      string     `json:"env_name" form:"env_name"`
	Docker       Docker     `json:"docker" form:"docker"`
	Repository   Repository `json:"repository" form:"repository"`
	BuildScript  string     `json:"build_script" form:"build_script"`
	K8SNamespace string     `json:"k8s_namespace" form:"k8s_namespace"`
}

type Repository struct {
	GitUrl      string `json:"git_url" form:"git_url"`
	GitBranch   string `json:"git_branch" form:"git_branch"`
	GitCommitId string `json:"git_commit_id" form:"git_commit_id"`
}
type Docker struct {
	BaseImage      string `json:"base_image" form:"base_image"`
	DockerfilePath string `json:"dockerfile_path" form:"dockerfile_path"`
	ImageName      string `json:"image_name" form:"image_name"`
	ImageTag       string `json:"image_tag" form:"image_tag"`
	Registry       struct {
		URL         string `json:"url"`
		Credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"credentials"`
	} `json:"registry"`
}
