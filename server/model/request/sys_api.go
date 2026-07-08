package request

import "github.com/tanganyu1114/heimdallr-reborn/model"

// api分页条件查询及排序结构体
type SearchApiParams struct {
	model.SysApi
	PageInfo
	OrderKey string `json:"orderKey"`
	Desc     bool   `json:"desc"`
}
