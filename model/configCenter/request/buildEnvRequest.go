package request

import "DYCLOUD/model/common/request"

type BuildEnvRequest struct {
	request.PageInfo
}

type DeleteBuildEnvByIds struct {
	Ids []int `json:"ids" form:"ids"`
}
