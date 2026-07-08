package response

import (
	"github.com/tanganyu1114/heimdallr-reborn/server/model"
)

type SysUserResponse struct {
	User model.SysUser `json:"user"`
}

type LoginResponse struct {
	User      model.SysUser `json:"user"`
	Token     string        `json:"token"`
	ExpiresAt int64         `json:"expiresAt"`
}

type APIKeyResponse struct {
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}
