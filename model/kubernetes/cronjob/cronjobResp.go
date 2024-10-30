package cronjob

import (
	"DYCLOUD/model/common/request"
	v1 "k8s.io/api/batch/v1"
)

type CronJobListResponse struct {
	Items *[]v1.CronJob `json:"items" form:"items"`
	Total int           `json:"total" form:"total"`
	request.PageInfo
}

type DescribeCronJobResponse struct {
	Items *v1.CronJob `json:"items" form:"items"`
}
