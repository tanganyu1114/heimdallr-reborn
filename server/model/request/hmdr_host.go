package request

import "gin-vue-admin/model"

type HmdrHostSearch struct {
	model.HmdrHost
	PageInfo
}
