package request

import (
	"DYCLOUD/model/common/request"
	v1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
)

type K8sVeleroTaskListResponse struct {
	Items *[]v1.Schedule `json:"items" form:"items"`
	Total int            `json:"total" form:"total"`
	request.PageInfo
}
type DescribeK8sVeleroTaskResponse struct {
	Items *v1.Schedule `json:"items" form:"items"`
}