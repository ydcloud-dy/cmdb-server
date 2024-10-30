package serviceAccount

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/core/v1"
)

type ServiceAccountListResponse struct {
	Items *[]v1.ServiceAccount `json:"items" form:"items"`
	Total int                  `json:"total" form:"total"`
	request.PageInfo
}

type DescribeServiceAccountResponse struct {
	Items *v1.ServiceAccount `json:"items" form:"items"`
}
