package namespaces

import "DYCLOUD/model/common/request"

type GetNamespaceListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	FieldSelector string `json:"fieldSelector" form:"fieldSelector"`
	Keyword       string `json:"keyword" form:"keyword"`
	request.PageInfo
}

func (r *GetNamespaceListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeNamespaceReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	NamespaceName string `json:"NamespaceName" form:"NamespaceName"`
}

func (r *DescribeNamespaceReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteNamespaceReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	NamespaceName string `json:"NamespaceName" form:"NamespaceName"`
}

func (r *DeleteNamespaceReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateNamespaceReq struct {
	ClusterId     int         `json:"cluster_id" form:"cluster_id"`
	Namespace     string      `json:"namespace" form:"namespace"`
	NamespaceName string      `json:"NamespaceName" form:"NamespaceName"`
	Content       interface{} `json:"content" form:"content"`
}

func (r *UpdateNamespaceReq) GetClusterID() int {
	return r.ClusterId
}

type CreateNamespaceReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateNamespaceReq) GetClusterID() int {
	return r.ClusterId
}
