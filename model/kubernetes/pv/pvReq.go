package pv

import "DYCLOUD/model/common/request"

type GetPVListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	request.PageInfo
}

func (r *GetPVListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribePVReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	PVName    string `json:"pvName" form:"pvName"`
}

func (r *DescribePVReq) GetClusterID() int {
	return r.ClusterId
}

type DeletePVReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	PVName    string `json:"pvName" form:"pvName"`
}

func (r *DeletePVReq) GetClusterID() int {
	return r.ClusterId
}

type UpdatePVReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	PVName    string      `json:"pvName" form:"pvName"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *UpdatePVReq) GetClusterID() int {
	return r.ClusterId
}

type CreatePVReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreatePVReq) GetClusterID() int {
	return r.ClusterId
}
