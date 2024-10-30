package storageClass

import "DYCLOUD/model/common/request"

type GetStorageClassListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	request.PageInfo
}

func (r *GetStorageClassListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeStorageClassReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`
	Namespace        string `json:"namespace" form:"namespace"`
	StorageClassName string `json:"storageClassName" form:"storageClassName"`
}

func (r *DescribeStorageClassReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteStorageClassReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`
	Namespace        string `json:"namespace" form:"namespace"`
	StorageClassName string `json:"storageClassName" form:"storageClassName"`
}

func (r *DeleteStorageClassReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateStorageClassReq struct {
	ClusterId        int         `json:"cluster_id" form:"cluster_id"`
	Namespace        string      `json:"namespace" form:"namespace"`
	StorageClassName string      `json:"storageClassName" form:"storageClassName"`
	Content          interface{} `json:"content" form:"content"`
}

func (r *UpdateStorageClassReq) GetClusterID() int {
	return r.ClusterId
}

type CreateStorageClassReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateStorageClassReq) GetClusterID() int {
	return r.ClusterId
}
