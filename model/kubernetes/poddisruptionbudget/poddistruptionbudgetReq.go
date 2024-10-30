package poddisruptionbudget

import "DYCLOUD/model/common/request"

type GetPoddisruptionbudgetListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`

	request.PageInfo
}

func (r *GetPoddisruptionbudgetListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribePoddisruptionbudgetReq struct {
	ClusterId               int    `json:"cluster_id" form:"cluster_id"`
	Namespace               string `json:"namespace" form:"namespace"`
	PoddisruptionbudgetName string `json:"poddisruptionbudgetName" form:"poddisruptionbudgetName"`
}

func (r *DescribePoddisruptionbudgetReq) GetClusterID() int {
	return r.ClusterId
}

type DeletePoddisruptionbudgetReq struct {
	ClusterId               int    `json:"cluster_id" form:"cluster_id"`
	Namespace               string `json:"namespace" form:"namespace"`
	PoddisruptionbudgetName string `json:"poddisruptionbudgetName" form:"poddisruptionbudgetName"`
}

func (r *DeletePoddisruptionbudgetReq) GetClusterID() int {
	return r.ClusterId
}

type UpdatePoddisruptionbudgetReq struct {
	ClusterId               int         `json:"cluster_id" form:"cluster_id"`
	Namespace               string      `json:"namespace" form:"namespace"`
	PoddisruptionbudgetName string      `json:"poddisruptionbudgetName" form:"poddisruptionbudgetName"`
	Content                 interface{} `json:"content" form:"content"`
}

func (r *UpdatePoddisruptionbudgetReq) GetClusterID() int {
	return r.ClusterId
}

type CreatePoddisruptionbudgetReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreatePoddisruptionbudgetReq) GetClusterID() int {
	return r.ClusterId
}
