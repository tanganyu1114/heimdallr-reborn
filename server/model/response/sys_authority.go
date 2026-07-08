package response

import "github.com/tanganyu1114/heimdallr-reborn/server/model"

type SysAuthorityResponse struct {
	Authority model.SysAuthority `json:"authority"`
}

type SysAuthorityCopyResponse struct {
	Authority      model.SysAuthority `json:"authority"`
	OldAuthorityId string             `json:"oldAuthorityId"`
}
