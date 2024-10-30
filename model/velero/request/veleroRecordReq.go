package request

import (
	"DYCLOUD/model/common/request"
)

type K8sVeleroRecordsSearchReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	request.PageInfo
}

func (r *K8sVeleroRecordsSearchReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeVeleroRecordReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`
	Namespace        string `json:"namespace" form:"namespace"`
	VeleroRecordName string `json:"VeleroRecordName" form:"VeleroRecordName"`
}

func (r *DescribeVeleroRecordReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteVeleroRecordReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`
	Namespace        string `json:"namespace" form:"namespace"`
	VeleroRecordName string `json:"VeleroRecordName" form:"VeleroRecordName"`
}

func (r *DeleteVeleroRecordReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateVeleroRecordReq struct {
	ClusterId        int         `json:"cluster_id" form:"cluster_id"`
	Namespace        string      `json:"namespace" form:"namespace"`
	VeleroRecordName string      `json:"VeleroRecordName" form:"VeleroRecordName"`
	Content          interface{} `json:"content" form:"content"`
}

func (r *UpdateVeleroRecordReq) GetClusterID() int {
	return r.ClusterId
}

type CreateVeleroRecordReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateVeleroRecordReq) GetClusterID() int {
	return r.ClusterId
}
