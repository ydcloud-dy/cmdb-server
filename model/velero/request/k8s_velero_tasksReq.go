package request

import (
	"DYCLOUD/model/common/request"
)

type K8sVeleroTasksSearchReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`

	request.PageInfo
}

func (r *K8sVeleroTasksSearchReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeVeleroTaskReq struct {
	ClusterId      int    `json:"cluster_id" form:"cluster_id"`
	Namespace      string `json:"namespace" form:"namespace"`
	VeleroTaskName string `json:"VeleroTaskName" form:"VeleroTaskName"`
}

func (r *DescribeVeleroTaskReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteVeleroTaskReq struct {
	ClusterId      int    `json:"cluster_id" form:"cluster_id"`
	Namespace      string `json:"namespace" form:"namespace"`
	VeleroTaskName string `json:"VeleroTaskName" form:"VeleroTaskName"`
}

func (r *DeleteVeleroTaskReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateVeleroTaskReq struct {
	ClusterId      int         `json:"cluster_id" form:"cluster_id"`
	Namespace      string      `json:"namespace" form:"namespace"`
	VeleroTaskName string      `json:"VeleroTaskName" form:"VeleroTaskName"`
	Content        interface{} `json:"content" form:"content"`
}

func (r *UpdateVeleroTaskReq) GetClusterID() int {
	return r.ClusterId
}

type CreateVeleroTaskReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateVeleroTaskReq) GetClusterID() int {
	return r.ClusterId
}
