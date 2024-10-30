package namespaces

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/core/v1"
)

type NamespaceListResponse struct {
	Items *[]v1.Namespace `json:"items" form:"items"`
	Total int             `json:"total" form:"total"`
	request.PageInfo
}

type DescribeNamespaceResponse struct {
	Items *v1.Namespace `json:"items" form:"items"`
}
