package request

import (
	"DYCLOUD/model/common/request"
	"time"
)

type CmdbHostsSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	Project        int        `json:"project" form:"project"`
	request.PageInfo
}
