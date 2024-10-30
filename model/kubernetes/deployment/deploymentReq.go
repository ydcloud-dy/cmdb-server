package deployment

import (
	"DYCLOUD/model/common/request"
)

type GetDeploymentListReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	request.PageInfo
}

func (r *GetDeploymentListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeDeploymentInfoReq struct {
	ClusterId      int    `json:"cluster_id" form:"cluster_id"`
	Namespace      string `json:"namespace" form:"namespace"`
	DeploymentName string `json:"deploymentName" form:"deploymentName"`
}

func (r *DescribeDeploymentInfoReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateDeploymentInfoReq struct {
	ClusterId      int         `json:"cluster_id" form:"cluster_id"`
	Namespace      string      `json:"namespace" form:"namespace"`
	DeploymentName string      `json:"deploymentName" form:"deploymentName"`
	Content        interface{} `json:"content" form:"content"`
}

func (r *UpdateDeploymentInfoReq) GetClusterID() int {
	return r.ClusterId
}

type CreateDeploymentReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreateDeploymentReq) GetClusterID() int {
	return r.ClusterId
}

type DeleteDeploymentReq struct {
	ClusterId      int    `json:"cluster_id" form:"cluster_id"`
	Namespace      string `json:"namespace" form:"namespace"`
	DeploymentName string `json:"deploymentName" form:"deploymentName"`
}

func (r *DeleteDeploymentReq) GetClusterID() int {
	return r.ClusterId
}

type RollBackDeployment struct {
	ClusterId      int         `json:"cluster_id" form:"cluster_id"`
	Namespace      string      `json:"namespace" form:"namespace"`
	DeploymentName string      `json:"deploymentName" form:"deploymentName"`
	Content        interface{} `json:"content" form:"content"`
}

func (r *RollBackDeployment) GetClusterID() int {
	return r.ClusterId
}
