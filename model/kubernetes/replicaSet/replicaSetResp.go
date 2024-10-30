package replicaSet

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/apps/v1"
)

type ReplicaSetListResponse struct {
	Items *[]v1.ReplicaSet `json:"items" form:"items"`
	Total int              `json:"total" form:"total"`
	request.PageInfo
}
