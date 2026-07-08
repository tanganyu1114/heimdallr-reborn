package response

import "github.com/tanganyu1114/heimdallr-reborn/model"

type SysAPIResponse struct {
	Api model.SysApi `json:"api"`
}

type SysAPIListResponse struct {
	Apis []model.SysApi `json:"apis"`
}
