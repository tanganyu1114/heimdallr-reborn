package request

import "github.com/tanganyu1114/heimdallr-reborn/model"

// Add menu authority info structure
type AddMenuAuthorityInfo struct {
	Menus       []model.SysBaseMenu
	AuthorityId string
}
