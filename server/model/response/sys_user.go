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

type SDKChallengeResponse struct {
	PublicKey string `json:"publicKey"` // RSA公钥
	Challenge string `json:"challenge"` // 挑战码
}

type APIKeyResponse struct {
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}
