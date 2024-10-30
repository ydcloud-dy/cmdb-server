package configmap

import "DYCLOUD/model/common/request"

type GetConfigMapListReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	request.PageInfo
}

func (r *GetConfigMapListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeConfigMapReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	ConfigMapName string `json:"configmapName" form:"configmapName"`
}

func (r *DescribeConfigMapReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteConfigMapReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	ConfigMapName string `json:"configmapName" form:"configmapName"`
}

func (r *DeleteConfigMapReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateConfigMapReq struct {
	ClusterId     int         `json:"cluster_id" form:"cluster_id"`
	Namespace     string      `json:"namespace" form:"namespace"`
	ConfigMapName string      `json:"configmapName" form:"configmapName"`
	Content       interface{} `json:"content" form:"content"`
}

func (r *UpdateConfigMapReq) GetClusterID() int {
	return r.ClusterId
}

type CreateConfigMapReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateConfigMapReq) GetClusterID() int {
	return r.ClusterId
}
