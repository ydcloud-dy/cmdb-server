package pods

import (
	"DYCLOUD/model/common/request"
	"DYCLOUD/utils/kubernetes/podtool"
	corev1 "k8s.io/api/core/v1"
)

type PodListResponse struct {
	Items *[]corev1.Pod `json:"items" form:"items"`
	Total int           `json:"total" form:"total"`
	request.PageInfo
}

type DescribePodInfoResponse struct {
	Items *corev1.Pod `json:"items" form:"items"`
}

type EventInfoResponse struct {
	Items *[]corev1.Event `json:"items" form:"items"`
	Total int             `json:"total" form:"total"`
}

type PodFilesResponse struct {
	Files []podtool.File `json:"files" form:"files"`
}
