package daemonset

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/apps/v1"
)

type DaemonSetListResponse struct {
	Items *[]v1.DaemonSet `json:"items" form:"items"`
	Total int             `json:"total" form:"total"`
	request.PageInfo
}

type DescribeDaemonSetResponse struct {
	Items *v1.DaemonSet `json:"items" form:"items"`
}
