package response

import "github.com/tanganyu1114/heimdallr-reborn/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
