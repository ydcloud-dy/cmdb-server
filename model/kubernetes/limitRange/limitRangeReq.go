package limitRange

import "DYCLOUD/model/common/request"

type GetLimitRangeListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	request.PageInfo
}

func (r *GetLimitRangeListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeLimitRangeReq struct {
	ClusterId      int    `json:"cluster_id" form:"cluster_id"`
	Namespace      string `json:"namespace" form:"namespace"`
	LimitRangeName string `json:"limitRangeName" form:"limitRangeName"`
}

func (r *DescribeLimitRangeReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteLimitRangeReq struct {
	ClusterId      int    `json:"cluster_id" form:"cluster_id"`
	Namespace      string `json:"namespace" form:"namespace"`
	LimitRangeName string `json:"limitRangeName" form:"limitRangeName"`
}

func (r *DeleteLimitRangeReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateLimitRangeReq struct {
	ClusterId      int         `json:"cluster_id" form:"cluster_id"`
	Namespace      string      `json:"namespace" form:"namespace"`
	LimitRangeName string      `json:"limitRangeName" form:"limitRangeName"`
	Content        interface{} `json:"content" form:"content"`
}

func (r *UpdateLimitRangeReq) GetClusterID() int {
	return r.ClusterId
}

type CreateLimitRangeReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateLimitRangeReq) GetClusterID() int {
	return r.ClusterId
}
