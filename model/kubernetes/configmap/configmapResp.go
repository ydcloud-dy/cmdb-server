package configmap

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/core/v1"
)

type ConfigMapListResponse struct {
	Items *[]v1.ConfigMap `json:"items" form:"items"`
	Total int             `json:"total" form:"total"`
	request.PageInfo
}
type DescribeConfigMapResponse struct {
	Items *v1.ConfigMap `json:"items"`
}
