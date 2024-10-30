package service

import "DYCLOUD/model/common/request"

type GetServiceListReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	request.PageInfo
}

func (r *GetServiceListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeServiceReq struct {
	ClusterId   int    `json:"cluster_id" form:"cluster_id"`
	Namespace   string `json:"namespace" form:"namespace"`
	ServiceName string `json:"serviceName" form:"serviceName"`
}

func (r *DescribeServiceReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteServiceReq struct {
	ClusterId   int    `json:"cluster_id" form:"cluster_id"`
	ServiceName string `json:"serviceName" form:"serviceName"`
	Namespace   string `json:"namespace" form:"namespace"`
}

func (r *DeleteServiceReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateServiceReq struct {
	ClusterId   int         `json:"cluster_id" form:"cluster_id"`
	Namespace   string      `json:"namespace" form:"namespace"`
	ServiceName string      `json:"serviceName" form:"serviceName"`
	Content     interface{} `json:"content" form:"content"`
}

func (r *UpdateServiceReq) GetClusterID() int {
	return r.ClusterId
}

type CreateServiceReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateServiceReq) GetClusterID() int {
	return r.ClusterId
}
