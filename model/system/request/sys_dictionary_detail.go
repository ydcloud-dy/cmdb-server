package request

import (
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
