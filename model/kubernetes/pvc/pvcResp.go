package pvc

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/core/v1"
)

type PvcListResponse struct {
	Items *[]v1.PersistentVolumeClaim `json:"items" form:"items"`
	Total int                         `json:"total" form:"total"`
	request.PageInfo
}
type DescribePVCResponse struct {
	Items *v1.PersistentVolumeClaim `json:"items" form:"items"`
}
