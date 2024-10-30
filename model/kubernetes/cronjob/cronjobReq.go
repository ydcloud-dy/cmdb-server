package cronjob

import "DYCLOUD/model/common/request"

type GetCronJobListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	FieldSelector string `json:"fieldSelector" form:"fieldSelector"`
	Keyword       string `json:"keyword" form:"keyword"`
	request.PageInfo
}

func (r *GetCronJobListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeCronJobReq struct {
	ClusterId   int    `json:"cluster_id" form:"cluster_id"`
	Namespace   string `json:"namespace" form:"namespace"`
	CronJobName string `json:"cronjobName" form:"cronjobName"`
}

func (r *DescribeCronJobReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteCronJobReq struct {
	ClusterId   int    `json:"cluster_id" form:"cluster_id"`
	Namespace   string `json:"namespace" form:"namespace"`
	CronJobName string `json:"cronjobName" form:"cronjobName"`
}

func (r *DeleteCronJobReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateCronJobReq struct {
	ClusterId   int         `json:"cluster_id" form:"cluster_id"`
	Namespace   string      `json:"namespace" form:"namespace"`
	CronJobName string      `json:"cronjobName" form:"cronjobName"`
	Content     interface{} `json:"content" form:"content"`
}

func (r *UpdateCronJobReq) GetClusterID() int {
	return r.ClusterId
}

type CreateCronJobReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateCronJobReq) GetClusterID() int {
	return r.ClusterId
}
