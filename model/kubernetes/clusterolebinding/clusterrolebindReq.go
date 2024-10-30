package clusterolebinding

import "DYCLOUD/model/common/request"

type GetClusterRoleBindingListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	request.PageInfo
}

func (r *GetClusterRoleBindingListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeClusterRoleBindingReq struct {
	ClusterId              int    `json:"cluster_id" form:"cluster_id"`
	Namespace              string `json:"namespace" form:"namespace"`
	ClusterRoleBindingName string `json:"clusterRoleBindingName" form:"clusterRoleBindingName"`
}

func (r *DescribeClusterRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteClusterRoleBindingReq struct {
	ClusterId              int    `json:"cluster_id" form:"cluster_id"`
	Namespace              string `json:"namespace" form:"namespace"`
	ClusterRoleBindingName string `json:"clusterRoleBindingName" form:"clusterRoleBindingName"`
}

func (r *DeleteClusterRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateClusterRoleBindingReq struct {
	ClusterId              int         `json:"cluster_id" form:"cluster_id"`
	Namespace              string      `json:"namespace" form:"namespace"`
	ClusterRoleBindingName string      `json:"clusterRoleBindingName" form:"clusterRoleBindingName"`
	Content                interface{} `json:"content" form:"content"`
}

func (r *UpdateClusterRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}

type CreateClusterRoleBindingReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateClusterRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}
