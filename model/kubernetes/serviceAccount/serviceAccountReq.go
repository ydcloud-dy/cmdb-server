package serviceAccount

import "DYCLOUD/model/common/request"

type GetServiceAccountReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	request.PageInfo
}

func (r *GetServiceAccountReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeServiceAccountReq struct {
	ClusterId          int    `json:"cluster_id" form:"cluster_id"`
	Namespace          string `json:"namespace" form:"namespace"`
	ServiceAccountName string `json:"serviceAccountName" form:"serviceAccountName"`
}

func (r *DescribeServiceAccountReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteServiceAccountReq struct {
	ClusterId          int    `json:"cluster_id" form:"cluster_id"`
	Namespace          string `json:"namespace" form:"namespace"`
	ServiceAccountName string `json:"serviceAccountName" form:"serviceAccountName"`
}

func (r *DeleteServiceAccountReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateServiceAccountReq struct {
	ClusterId          int         `json:"cluster_id" form:"cluster_id"`
	Namespace          string      `json:"namespace" form:"namespace"`
	ServiceAccountName string      `json:"serviceAccountName" form:"serviceAccountName"`
	Content            interface{} `json:"content" form:"content"`
}

func (r *UpdateServiceAccountReq) GetClusterID() int {
	return r.ClusterId
}

type CreateServiceAccountReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateServiceAccountReq) GetClusterID() int {
	return r.ClusterId
}
