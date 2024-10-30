package roles

import "DYCLOUD/model/common/request"

type GetRolesListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	request.PageInfo
}

func (r *GetRolesListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeRolesReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	RolesName string `json:"roleName" form:"roleName"`
}

func (r *DescribeRolesReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteRolesReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	RolesName string `json:"roleName" form:"roleName"`
}

func (r *DeleteRolesReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateRolesReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	RolesName string      `json:"roleName" form:"roleName"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *UpdateRolesReq) GetClusterID() int {
	return r.ClusterId
}

type CreateRolesReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateRolesReq) GetClusterID() int {
	return r.ClusterId
}
