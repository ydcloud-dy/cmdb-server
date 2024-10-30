package poddisruptionbudget

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/policy/v1"
)

type PoddisruptionbudgetListResponse struct {
	Items *[]v1.PodDisruptionBudget `json:"items" form:"items"`
	Total int                       `json:"total" form:"total"`
	request.PageInfo
}
type DescribePoddisruptionbudgetResponse struct {
	Items *v1.PodDisruptionBudget `json:"items" form:"items"`
}
