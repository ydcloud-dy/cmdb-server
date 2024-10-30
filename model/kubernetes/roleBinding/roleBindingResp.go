package roleBinding

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/rbac/v1"
)

type RoleBindingListResponse struct {
	Items *[]v1.RoleBinding `json:"items" form:"items"`
	Total int               `json:"total" form:"total"`
	request.PageInfo
}
type DescribeRoleBindingResponse struct {
	Items *v1.RoleBinding `json:"items" form:"items"`
}
