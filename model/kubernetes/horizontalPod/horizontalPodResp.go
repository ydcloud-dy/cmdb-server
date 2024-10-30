package horizontalPod

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/autoscaling/v1"
)

type HorizontalPodResponse struct {
	Items *[]v1.HorizontalPodAutoscaler `json:"items" form:"items"`
	Total int                           `json:"total" form:"total"`
	request.PageInfo
}
type DescribeHorizontalPodResponse struct {
	Items *v1.HorizontalPodAutoscaler `json:"items" form:"items"`
}
