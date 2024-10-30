package horizontalPod

import "DYCLOUD/model/common/request"

type GetHorizontalPodListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	request.PageInfo
}

func (r *GetHorizontalPodListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeHorizontalPodReq struct {
	ClusterId         int    `json:"cluster_id" form:"cluster_id"`
	Namespace         string `json:"namespace" form:"namespace"`
	HorizontalPodName string `json:"HorizontalPodName" form:"HorizontalPodName"`
}

func (r *DescribeHorizontalPodReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteHorizontalPodReq struct {
	ClusterId         int    `json:"cluster_id" form:"cluster_id"`
	Namespace         string `json:"namespace" form:"namespace"`
	HorizontalPodName string `json:"HorizontalPodName" form:"HorizontalPodName"`
}

func (r *DeleteHorizontalPodReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateHorizontalPodReq struct {
	ClusterId         int         `json:"cluster_id" form:"cluster_id"`
	Namespace         string      `json:"namespace" form:"namespace"`
	HorizontalPodName string      `json:"HorizontalPodName" form:"HorizontalPodName"`
	Content           interface{} `json:"content" form:"content"`
}

func (r *UpdateHorizontalPodReq) GetClusterID() int {
	return r.ClusterId
}

type CreateHorizontalPodReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateHorizontalPodReq) GetClusterID() int {
	return r.ClusterId
}
