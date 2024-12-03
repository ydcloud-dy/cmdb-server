package request

import (
	"DYCLOUD/model/common/request"
	"time"
)

type PipelinesRequest struct {
	Namespace  string `json:"namespace" form:"namespace"`
	Cluster_ID int    `json:"cluster_id" form:"cluster_id" gorm:"-"`
	AppCode    string `json:"appCode" form:"appCode" gorm:"-"`
	EnvCode    string `json:"envCode" form:"envCode" gorm:"-"`
	request.PageInfo
	StartCreatedAt time.Time
	EndCreatedAt   time.Time
}
type PipelineRunStatus struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	User        string `json:"user"`
	Branch      string `json:"branch"`
	LastRunTime string `json:"last_run_time"`
	Duration    string `json:"duration"`
}
