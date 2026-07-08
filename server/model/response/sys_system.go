package response

import "github.com/tanganyu1114/heimdallr-reborn/server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
