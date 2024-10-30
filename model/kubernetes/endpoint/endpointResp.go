package endpoint

import (
	"DYCLOUD/model/common/request"
	corev1 "k8s.io/api/core/v1"
)

type EndPointListResponse struct {
	Items *[]corev1.Endpoints `json:"items" form:"items"`
	Total int                 `json:"total" form:"total"`
	request.PageInfo
}

type DescribeEndPointResponse struct {
	Items *corev1.Endpoints `json:"items" form:"items"`
}
