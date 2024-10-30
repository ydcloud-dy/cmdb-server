package clusterolebinding

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/rbac/v1"
)

type ClusterRoleBindingListResponse struct {
	Items *[]v1.ClusterRoleBinding `json:"items" form:"items"`
	Total int                      `json:"total" form:"total"`
	request.PageInfo
}
type DescribeClusterRoleBindingResponse struct {
	Items *v1.ClusterRoleBinding `json:"items" form:"items"`
}
