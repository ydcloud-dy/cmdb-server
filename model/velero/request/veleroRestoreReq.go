package request

import (
	"DYCLOUD/model/common/request"
)

type K8sVeleroRestoresSearchReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	request.PageInfo
}

func (r *K8sVeleroRestoresSearchReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeVeleroRestoreReq struct {
	ClusterId         int    `json:"cluster_id" form:"cluster_id"`
	Namespace         string `json:"namespace" form:"namespace"`
	VeleroRestoreName string `json:"VeleroRestoreName" form:"VeleroRestoreName"`
}

func (r *DescribeVeleroRestoreReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteVeleroRestoreReq struct {
	ClusterId         int    `json:"cluster_id" form:"cluster_id"`
	Namespace         string `json:"namespace" form:"namespace"`
	VeleroRestoreName string `json:"VeleroRestoreName" form:"VeleroRestoreName"`
}

func (r *DeleteVeleroRestoreReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateVeleroRestoreReq struct {
	ClusterId         int         `json:"cluster_id" form:"cluster_id"`
	Namespace         string      `json:"namespace" form:"namespace"`
	VeleroRestoreName string      `json:"VeleroRestoreName" form:"VeleroRestoreName"`
	Content           interface{} `json:"content" form:"content"`
}

func (r *UpdateVeleroRestoreReq) GetClusterID() int {
	return r.ClusterId
}

type CreateVeleroRestoreReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateVeleroRestoreReq) GetClusterID() int {
	return r.ClusterId
}
