package pvc

import "DYCLOUD/model/common/request"

type GetPvcListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	request.PageInfo
}

func (r *GetPvcListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribePVCReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	PVCName   string `json:"pvcName" form:"pvcName"`
}

func (r *DescribePVCReq) GetClusterID() int {
	return r.ClusterId
}

type DeletePVCReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	PVCName   string `json:"pvcName" form:"pvcName"`
}

func (r *DeletePVCReq) GetClusterID() int {
	return r.ClusterId
}

type UpdatePVCReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	PVCName   string      `json:"pvcName" form:"pvcName"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *UpdatePVCReq) GetClusterID() int {
	return r.ClusterId
}

type CreatePVCReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreatePVCReq) GetClusterID() int {
	return r.ClusterId
}
