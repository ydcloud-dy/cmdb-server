package secret

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/core/v1"
)

type SecretListResponse struct {
	Items *[]v1.Secret `json:"items" form:"items"`
	Total int          `json:"total" form:"total"`
	request.PageInfo
}
type DescribeSecretResponse struct {
	Items *v1.Secret `json:"items"`
}
