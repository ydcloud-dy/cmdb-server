package replicaSet

import "DYCLOUD/model/common/request"

type GetReplicaSetListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	request.PageInfo
}

func (r *GetReplicaSetListReq) GetClusterID() int {
	return r.ClusterId
}
