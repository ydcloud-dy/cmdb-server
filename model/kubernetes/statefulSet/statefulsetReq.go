package statefulSet

import "DYCLOUD/model/common/request"

type GetStatefulSetListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	FieldSelector string `json:"fieldSelector" form:"fieldSelector"`
	Keyword       string `json:"keyword" form:"keyword"`
	request.PageInfo
}

func (r *GetStatefulSetListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeStatefulSetReq struct {
	ClusterId       int    `json:"cluster_id" form:"cluster_id"`
	Namespace       string `json:"namespace" form:"namespace"`
	StatefulSetName string `json:"statefulsetName" form:"statefulsetName"`
}

func (r *DescribeStatefulSetReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteStatefulSetReq struct {
	ClusterId       int    `json:"cluster_id" form:"cluster_id"`
	Namespace       string `json:"namespace" form:"namespace"`
	StatefulSetName string `json:"statefulsetName" form:"statefulsetName"`
}

func (r *DeleteStatefulSetReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateStatefulSetReq struct {
	ClusterId       int         `json:"cluster_id" form:"cluster_id"`
	Namespace       string      `json:"namespace" form:"namespace"`
	StatefulSetName string      `json:"statefulsetName" form:"statefulsetName"`
	Content         interface{} `json:"content" form:"content"`
}

func (r *UpdateStatefulSetReq) GetClusterID() int {
	return r.ClusterId
}

type CreateStatefulSetReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateStatefulSetReq) GetClusterID() int {
	return r.ClusterId
}
