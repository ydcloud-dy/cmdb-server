package secret

import "DYCLOUD/model/common/request"

type GetSecretList struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	request.PageInfo
}

func (r *GetSecretList) GetClusterID() int {
	return r.ClusterId
}

type DescribeSecretReq struct {
	ClusterId  int    `json:"cluster_id" form:"cluster_id"`
	Namespace  string `json:"namespace" form:"namespace"`
	SecretName string `json:"secretName" form:"secretName"`
}

func (r *DescribeSecretReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteSecretReq struct {
	ClusterId  int    `json:"cluster_id" form:"cluster_id"`
	Namespace  string `json:"namespace" form:"namespace"`
	SecretName string `json:"secretName" form:"secretName"`
}

func (r *DeleteSecretReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateSecretReq struct {
	ClusterId  int         `json:"cluster_id" form:"cluster_id"`
	Namespace  string      `json:"namespace" form:"namespace"`
	SecretName string      `json:"secretName" form:"secretName"`
	Content    interface{} `json:"content" form:"content"`
}

func (r *UpdateSecretReq) GetClusterID() int {
	return r.ClusterId
}

type CreateSecretReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateSecretReq) GetClusterID() int {
	return r.ClusterId
}
