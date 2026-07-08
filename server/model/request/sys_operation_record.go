package request

import "github.com/tanganyu1114/heimdallr-reborn/server/model"

type SysOperationRecordSearch struct {
	model.SysOperationRecord
	PageInfo
}
