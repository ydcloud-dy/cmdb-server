package clusterrole

import "DYCLOUD/model/common/request"

type GetClusterRoleListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	request.PageInfo
}

func (r *GetClusterRoleListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeClusterRoleReq struct {
	ClusterId       int    `json:"cluster_id" form:"cluster_id"`
	Namespace       string `json:"namespace" form:"namespace"`
	ClusterRoleName string `json:"clusterRoleName" form:"clusterRoleName"`
}

func (r *DescribeClusterRoleReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteClusterRoleReq struct {
	ClusterId       int    `json:"cluster_id" form:"cluster_id"`
	Namespace       string `json:"namespace" form:"namespace"`
	ClusterRoleName string `json:"clusterRoleName" form:"clusterRoleName"`
}

func (r *DeleteClusterRoleReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateClusterRoleReq struct {
	ClusterId       int         `json:"cluster_id" form:"cluster_id"`
	Namespace       string      `json:"namespace" form:"namespace"`
	ClusterRoleName string      `json:"clusterRoleName" form:"clusterRoleName"`
	Content         interface{} `json:"content" form:"content"`
}

func (r *UpdateClusterRoleReq) GetClusterID() int {
	return r.ClusterId
}

type CreateClusterRoleReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateClusterRoleReq) GetClusterID() int {
	return r.ClusterId
}
