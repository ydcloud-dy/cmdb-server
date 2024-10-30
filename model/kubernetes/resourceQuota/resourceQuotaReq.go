package resourceQuota

import "DYCLOUD/model/common/request"

type GetResourceQuotaListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	request.PageInfo
}

func (r *GetResourceQuotaListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeResourceQuotaReq struct {
	ClusterId         int    `json:"cluster_id" form:"cluster_id"`
	Namespace         string `json:"namespace" form:"namespace"`
	ResourceQuotaName string `json:"resourcequotaName" form:"resourcequotaName"`
}

func (r *DescribeResourceQuotaReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteResourceQuotaReq struct {
	ClusterId         int    `json:"cluster_id" form:"cluster_id"`
	Namespace         string `json:"namespace" form:"namespace"`
	ResourceQuotaName string `json:"resourcequotaName" form:"resourcequotaName"`
}

func (r *DeleteResourceQuotaReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateResourceQuotaReq struct {
	ClusterId         int         `json:"cluster_id" form:"cluster_id"`
	Namespace         string      `json:"namespace" form:"namespace"`
	ResourceQuotaName string      `json:"resourcequotaName" form:"resourcequotaName"`
	Content           interface{} `json:"content" form:"content"`
}

func (r *UpdateResourceQuotaReq) GetClusterID() int {
	return r.ClusterId
}

type CreateResourceQuotaReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateResourceQuotaReq) GetClusterID() int {
	return r.ClusterId
}
