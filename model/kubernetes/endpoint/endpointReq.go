package endpoint

import "DYCLOUD/model/common/request"

type GetEndPointListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	FieldSelector string `json:"fieldSelector" form:"fieldSelector"`
	Keyword       string `json:"keyword" form:"keyword"`
	request.PageInfo
}

func (r *GetEndPointListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeEndPointReq struct {
	ClusterId    int    `json:"cluster_id" form:"cluster_id"`
	Namespace    string `json:"namespace" form:"namespace"`
	EndPointName string `json:"endpointName" form:"endpointName"`
}

func (r *DescribeEndPointReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteEndPointReq struct {
	ClusterId    int    `json:"cluster_id" form:"cluster_id"`
	Namespace    string `json:"namespace" form:"namespace"`
	EndPointName string `json:"endpointName" form:"endpointName"`
}

func (r *DeleteEndPointReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateEndPointReq struct {
	ClusterId    int         `json:"cluster_id" form:"cluster_id"`
	Namespace    string      `json:"namespace" form:"namespace"`
	EndPointName string      `json:"endpointName" form:"endpointName"`
	Content      interface{} `json:"content" form:"content"`
}

func (r *UpdateEndPointReq) GetClusterID() int {
	return r.ClusterId
}

type CreateEndPointReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateEndPointReq) GetClusterID() int {
	return r.ClusterId
}
