package daemonset

import "DYCLOUD/model/common/request"

type GetDaemonSetListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	FieldSelector string `json:"fieldSelector" form:"fieldSelector"`
	Keyword       string `json:"keyword" form:"keyword"`
	request.PageInfo
}

func (r *GetDaemonSetListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeDaemonSetReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	DaemonsetName string `json:"daemonsetName" form:"daemonsetName"`
}

func (r *DescribeDaemonSetReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteDaemonSetReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	DaemonsetName string `json:"daemonsetName" form:"daemonsetName"`
}

func (r *DeleteDaemonSetReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateDaemonSetReq struct {
	ClusterId     int         `json:"cluster_id" form:"cluster_id"`
	Namespace     string      `json:"namespace" form:"namespace"`
	DaemonsetName string      `json:"daemonsetName" form:"daemonsetName"`
	Content       interface{} `json:"content" form:"content"`
}

func (r *UpdateDaemonSetReq) GetClusterID() int {
	return r.ClusterId
}

type CreateDaemonSetReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateDaemonSetReq) GetClusterID() int {
	return r.ClusterId
}
