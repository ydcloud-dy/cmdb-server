package service

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/core/v1"
)

type ServiceListResponse struct {
	Items *[]v1.Service `json:"items" form:"items"`
	Total int           `json:"total" form:"total"`
	request.PageInfo
}

type DescribeServiceResponse struct {
	Items *v1.Service `json:"items" form:"items"`
}
