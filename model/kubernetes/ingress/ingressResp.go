package ingress

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/networking/v1"
)

type IngressListResponse struct {
	Items *[]v1.Ingress `json:"items" form:"items"`
	Total int           `json:"total" form:"total"`
	request.PageInfo
}
type DescribeIngressResponse struct {
	Items *v1.Ingress `json:"items" form:"items"`
}
