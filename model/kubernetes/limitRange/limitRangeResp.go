package limitRange

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/core/v1"
)

type LimitRangeListResponse struct {
	Items *[]v1.LimitRange `json:"items" form:"items"`
	Total int              `json:"total" form:"total"`
	request.PageInfo
}
type DescribeLimitRangeResponse struct {
	Items *v1.LimitRange `json:"items" form:"items"`
}
