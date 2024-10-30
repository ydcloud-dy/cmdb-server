package clusterrole

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/rbac/v1"
)

type ClusterRoleListResponse struct {
	Items *[]v1.ClusterRole `json:"items" form:"items"`
	Total int               `json:"total" form:"total"`
	request.PageInfo
}
type DescribeClusterRoleResponse struct {
	Items *v1.ClusterRole `json:"items" form:"items"`
}
