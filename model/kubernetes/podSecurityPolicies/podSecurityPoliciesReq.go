package podSecurityPolicies

import "DYCLOUD/model/common/request"

type GetPodSecurityPoliciesListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	request.PageInfo
}

func (r *GetPodSecurityPoliciesListReq) GetClusterID() int {
	return r.ClusterId
}

type DescribePodSecurityPoliciesReq struct {
	ClusterId               int    `json:"cluster_id" form:"cluster_id"`
	Namespace               string `json:"namespace" form:"namespace"`
	PodSecurityPoliciesName string `json:"podSecurityPoliciesName" form:"podSecurityPoliciesName"`
}

func (r *DescribePodSecurityPoliciesReq) GetClusterID() int {
	return r.ClusterId
}

type DeletePodSecurityPoliciesReq struct {
	ClusterId               int    `json:"cluster_id" form:"cluster_id"`
	Namespace               string `json:"namespace" form:"namespace"`
	PodSecurityPoliciesName string `json:"podSecurityPoliciesName" form:"podSecurityPoliciesName"`
}

func (r *DeletePodSecurityPoliciesReq) GetClusterID() int {
	return r.ClusterId
}

type UpdatePodSecurityPoliciesReq struct {
	ClusterId               int         `json:"cluster_id" form:"cluster_id"`
	Namespace               string      `json:"namespace" form:"namespace"`
	PodSecurityPoliciesName string      `json:"podSecurityPoliciesName" form:"podSecurityPoliciesName"`
	Content                 interface{} `json:"content" form:"content"`
}

func (r *UpdatePodSecurityPoliciesReq) GetClusterID() int {
	return r.ClusterId
}

type CreatePodSecurityPoliciesReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *CreatePodSecurityPoliciesReq) GetClusterID() int {
	return r.ClusterId
}
