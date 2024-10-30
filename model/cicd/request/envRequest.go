package request

import "DYCLOUD/model/common/request"

type EnvRequest struct {
	request.PageInfo
}

type DeleteEnvByIds struct {
	Ids []int `json:"ids" form:"ids"`
}
