package roleBinding

import "DYCLOUD/model/common/request"

type GetRoleBindingListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	request.PageInfo
}

func (r *GetRoleBindingListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeRoleBindingReq struct {
	ClusterId       int    `json:"cluster_id" form:"cluster_id"`
	Namespace       string `json:"namespace" form:"namespace"`
	RoleBindingName string `json:"roleBindingName" form:"roleBindingName"`
}

func (r *DescribeRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteRoleBindingReq struct {
	ClusterId       int    `json:"cluster_id" form:"cluster_id"`
	Namespace       string `json:"namespace" form:"namespace"`
	RoleBindingName string `json:"roleBindingName" form:"roleBindingName"`
}

func (r *DeleteRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateRoleBindingReq struct {
	ClusterId       int         `json:"cluster_id" form:"cluster_id"`
	Namespace       string      `json:"namespace" form:"namespace"`
	RoleBindingName string      `json:"roleBindingName" form:"roleBindingName"`
	Content         interface{} `json:"content" form:"content"`
}

func (r *UpdateRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}

type CreateRoleBindingReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}
