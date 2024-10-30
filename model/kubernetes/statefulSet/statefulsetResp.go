package statefulSet

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/apps/v1"
)

type StatefulSetListResponse struct {
	Items *[]v1.StatefulSet `json:"items" form:"items"`
	Total int               `json:"total" form:"total"`
	request.PageInfo
}

type DescribeStatefulSetResponse struct {
	Items *v1.StatefulSet `json:"items" form:"items"`
}
