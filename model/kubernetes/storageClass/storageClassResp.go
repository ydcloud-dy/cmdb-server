package storageClass

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/storage/v1"
)

type StorageClassListResponse struct {
	Items *[]v1.StorageClass `json:"items" form:"items"`
	Total int                `json:"total" form:"total"`
	request.PageInfo
}
type DescribeStorageClassResponse struct {
	Items *v1.StorageClass `json:"items" form:"items"`
}
