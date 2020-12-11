package request

import "gin-vue-admin/model"

type HmdrGroupSearch struct {
	model.HmdrGroup
	PageInfo
}
