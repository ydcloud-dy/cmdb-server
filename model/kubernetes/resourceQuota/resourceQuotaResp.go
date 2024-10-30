package resourceQuota

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/core/v1"
)

type ResourceQuotaListResponse struct {
	Items *[]v1.ResourceQuota `json:"items" form:"items"`
	Total int                 `json:"total" form:"total"`
	request.PageInfo
}
type DescribeResourceQuotaResponse struct {
	Items *v1.ResourceQuota `json:"items" form:"items"`
}
