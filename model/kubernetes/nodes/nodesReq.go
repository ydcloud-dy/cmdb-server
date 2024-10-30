package nodes

import (
	"DYCLOUD/model/common/request"
)

type NodeListReq struct {
	ClusterId int `json:"cluster_id" form:"cluster_id"`
	request.PageInfo
}

func (r *NodeListReq) GetClusterID() int {
	return r.ClusterId
}

type NodeMetricsReq struct {
	ClusterId int `json:"cluster_id" form:"cluster_id"`
}

func (r *NodeMetricsReq) GetClusterID() int {
	return r.ClusterId
}

type DescribeNodeReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	NodeName  string `json:"nodeName" form:"nodeName"`
}

func (r *DescribeNodeReq) GetClusterID() int {
	return r.ClusterId
}

type UpdateNodeReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	NodeName  string      `json:"nodeName" form:"nodeName"`
	Content   interface{} `json:"content" form:"content"`
}

func (r *UpdateNodeReq) GetClusterID() int {
	return r.ClusterId
}

type EvictAllNodePodReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	NodeName      string `json:"nodeName" form:"nodeName"`
	FieldSelector string `json:"fieldSelector" form:"fieldSelector"`
}

func (r *EvictAllNodePodReq) GetClusterID() int {
	return r.ClusterId
}
