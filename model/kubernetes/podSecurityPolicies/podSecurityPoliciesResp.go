package podSecurityPolicies

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/core/v1"
)

type PodSecurityPoliciesListResponse struct {
	Items *[]v1.PodSecurityContext `json:"items" form:"items"`
	Total int                      `json:"total" form:"total"`
	request.PageInfo
}

//type DescribeStorageClassResponse struct {
//	Items *v1.StorageClass `json:"items" form:"items"`
//}
