package roles

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/rbac/v1"
)

type RoleListResponse struct {
	Items *[]v1.Role `json:"items" form:"items"`
	Total int        `json:"total" form:"total"`
	request.PageInfo
}
type DescribeRoleResponse struct {
	Items *v1.Role `json:"items" form:"items"`
}
