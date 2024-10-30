package request

import "DYCLOUD/model/common/request"

type GetVeleroListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	FieldSelector string `json:"fieldSelector" form:"fieldSelector"`
	Keyword       string `json:"keyword" form:"keyword"`
	request.PageInfo
}

func (r *GetVeleroListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeVeleroReq struct {
	ClusterId  int    `json:"cluster_id" form:"cluster_id"`
	Namespace  string `json:"namespace" form:"namespace"`
	VeleroName string `json:"VeleroName" form:"VeleroName"`
}

func (r *DescribeVeleroReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteVeleroReq struct {
	ClusterId  int    `json:"cluster_id" form:"cluster_id"`
	Namespace  string `json:"namespace" form:"namespace"`
	VeleroName string `json:"VeleroName" form:"VeleroName"`
}

func (r *DeleteVeleroReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateVeleroReq struct {
	ClusterId  int         `json:"cluster_id" form:"cluster_id"`
	Namespace  string      `json:"namespace" form:"namespace"`
	VeleroName string      `json:"VeleroName" form:"VeleroName"`
	Content    interface{} `json:"content" form:"content"`
}

func (r *UpdateVeleroReq) GetClusterID() int {
	return r.ClusterId
}

type VeleroModel struct {
	ClusterId   int    `json:"cluster_id" form:"cluster_id"`
	S3Region    string `json:"s3Region" form:"s3Region"`
	S3Address   string `json:"s3Address" form:"s3Address"`
	S3Key       string `json:"s3Key" form:"s3Key"`
	S3Secret    string `json:"s3Secret" form:"s3Secret"`
	S3Bucket    string `json:"s3Bucket" form:"s3Bucket"`
	Provider    string `json:"provider"`
	VeleroImage string `json:"veleroImage" form:"veleroImage"`
	PluginImage string `json:"pluginImage" form:"pulginImage"`
}

func (v *VeleroModel) TableName() string {
	return "k8s_velero"
}
func (r *VeleroModel) GetClusterID() int {
	return r.ClusterId
}
