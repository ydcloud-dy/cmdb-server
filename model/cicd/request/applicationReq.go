package request

import (
	"DYCLOUD/model/common/request"
	"time"
)

type ApplicationRequest struct {
	AppId int `json:"app_id" form:"app_id" gorm:"-"`
	request.PageInfo
	StartCreatedAt time.Time
	EndCreatedAt   time.Time
}
type DeleteApplicationByIds struct {
	Ids []int `json:"ids" form:"ids"`
}
