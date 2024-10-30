package nodes

import (
	"DYCLOUD/model/common/request"
	corev1 "k8s.io/api/core/v1"
)

type NodeListResponse struct {
	Items *[]corev1.Node `json:"items" form:"items"`
	Total int            `json:"total" form:"total"`
	request.PageInfo
}

type DescribeNodeInfoResponse struct {
	Items *corev1.Node `json:"items" form:"items"`
}
