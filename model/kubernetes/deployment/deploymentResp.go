package deployment

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/apps/v1"
)

type DeploymentListResponse struct {
	Items *[]v1.Deployment `json:"items" form:"items"`
	Total int              `json:"total" form:"total"`
	request.PageInfo
}
type DescribeDeployResponse struct {
	Items *v1.Deployment `json:"items"`
}
