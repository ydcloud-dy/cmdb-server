package request

import (
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
