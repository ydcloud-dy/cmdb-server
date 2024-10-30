package pv

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/core/v1"
)

type PVListResponse struct {
	Items *[]v1.PersistentVolume `json:"items" form:"items"`
	Total int                    `json:"total" form:"total"`
	request.PageInfo
}
type DescribePVResponse struct {
	Items *v1.PersistentVolume `json:"items" form:"items"`
}
