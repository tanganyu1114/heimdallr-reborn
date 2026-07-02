package request

import uuid "github.com/satori/go.uuid"

// User register structure
type Register struct {
	Username    string `json:"userName"`
	Password    string `json:"passWord"`
	NickName    string `json:"nickName" gorm:"default:'QMPlusUser'"`
	HeaderImg   string `json:"headerImg" gorm:"default:'http://www.henrongyi.top/avatar/lufu.jpg'"`
	AuthorityId string `json:"authorityId" gorm:"default:888"`
}

// User login structure
type Login struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
}

// Modify password structure
type ChangePasswordStruct struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

// Modify  user's auth structure
type SetUserAuth struct {
	UUID        uuid.UUID `json:"uuid"`
	AuthorityId string    `json:"authorityId"`
}

// SDKLogin SDK登录请求结构（无需验证码）
type SDKLogin struct {
	APIKey    string `json:"apiKey"`    // API密钥
	APISecret string `json:"apiSecret"` // API密钥密码
}

// GenerateAPIKeyRequest 生成API Key请求
type GenerateAPIKeyRequest struct {
	UserID uint `json:"userId" binding:"required"` // 用户ID
}

// ToggleAPIKeyRequest 切换API Key状态请求
type ToggleAPIKeyRequest struct {
	UserID  uint `json:"userId" binding:"required"`  // 用户ID
	Enabled bool `json:"enabled" binding:"required"` // 是否启用
}

// RegenerateAPISecretRequest 重新生成API Secret请求
type RegenerateAPISecretRequest struct {
	UserID uint `json:"userId" binding:"required"` // 用户ID
}
