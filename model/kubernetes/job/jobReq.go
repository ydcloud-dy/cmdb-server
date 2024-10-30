package job

import "DYCLOUD/model/common/request"

type GetJobListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	FieldSelector string `json:"fieldSelector" form:"fieldSelector"`
	Keyword       string `json:"keyword" form:"keyword"`
	request.PageInfo
}

func (r *GetJobListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeJobReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	JobName   string `json:"jobName" form:"jobName"`
}

func (r *DescribeJobReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteJobReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	JobName   string `json:"jobName" form:"jobName"`
}

func (r *DeleteJobReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateJobReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	JobName   string      `json:"jobName" form:"jobName"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *UpdateJobReq) GetClusterID() int {
	return r.ClusterId
}

type CreateJobReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateJobReq) GetClusterID() int {
	return r.ClusterId
}
