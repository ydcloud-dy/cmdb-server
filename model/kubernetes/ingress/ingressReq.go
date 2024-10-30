package ingress

import "DYCLOUD/model/common/request"

type GetIngressListReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	request.PageInfo
}

func (r *GetIngressListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeIngressReq struct {
	ClusterId   int    `json:"cluster_id" form:"cluster_id"`
	Namespace   string `json:"namespace" form:"namespace"`
	IngressName string `json:"ingressName" form:"ingressName"`
}

func (r *DescribeIngressReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteIngressReq struct {
	ClusterId   int    `json:"cluster_id" form:"cluster_id"`
	Namespace   string `json:"namespace" form:"namespace"`
	IngressName string `json:"ingressName" form:"ingressName"`
}

func (r *DeleteIngressReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateIngressReq struct {
	ClusterId   int         `json:"cluster_id" form:"cluster_id"`
	Namespace   string      `json:"namespace" form:"namespace"`
	IngressName string      `json:"ingressName" form:"ingressName"`
	Content     interface{} `json:"content" form:"content"`
}

func (r *UpdateIngressReq) GetClusterID() int {
	return r.ClusterId
}

type CreateIngressReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateIngressReq) GetClusterID() int {
	return r.ClusterId
}
