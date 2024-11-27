package request

import (
	"DYCLOUD/model/common/request"
	"time"
)

type ApplicationRequest struct {
	AppId int `json:"app_id" form:"app_id" gorm:"-"`
	request.PageInfo
	StartCreatedAt time.Time
	EndCreatedAt   time.Time
}
type DeleteApplicationByIds struct {
	Ids []int `json:"ids" form:"ids"`
}

type DeploymentInfoRequest struct {
	ClusterId   int    `json:"cluster_id" form:"cluster_id"`
	Namespace   string `json:"namespace" form:"namespace"`
	ClusterName string `json:"cluster_name" form:"cluster_name"`
	EnvCode     string `json:"env_code" form:"env_code"`
	AppCode     string `json:"app_code" form:"app_code"`
}
