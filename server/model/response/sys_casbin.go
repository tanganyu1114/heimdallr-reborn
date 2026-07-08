package response

import "github.com/tanganyu1114/heimdallr-reborn/server/model/request"

type PolicyPathResponse struct {
	Paths []request.CasbinInfo `json:"paths"`
}
