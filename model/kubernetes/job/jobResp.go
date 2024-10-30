package job

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/batch/v1"
)

type JobListResponse struct {
	Items *[]v1.Job `json:"items" form:"items"`
	Total int       `json:"total" form:"total"`
	request.PageInfo
}

type DescribeJobResponse struct {
	Items *v1.Job `json:"items" form:"items"`
}
