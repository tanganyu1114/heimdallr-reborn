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
	Challenge string `json:"challenge"` // 挑战码，用于防重放攻击
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
	Challenge string `json:"challenge"` // 挑战码，用于防重放攻击
}

// SDKChallengeRequest 获取SDK挑战码请求
type SDKChallengeRequest struct {
	APIKey string `json:"apiKey" binding:"required"` // API密钥
}

// SDKChallengeResponse 获取SDK挑战码响应
type SDKChallengeResponse struct {
	PublicKey string `json:"publicKey"` // RSA公钥
	Challenge string `json:"challenge"` // 挑战码
}

// GetPublicKeyRequest 获取公钥和挑战码请求
type GetPublicKeyRequest struct {
	CaptchaId string `json:"captchaId" binding:"required"` // 验证码ID，用作会话标识
}

// GetPublicKeyResponse 获取公钥和挑战码响应
type GetPublicKeyResponse struct {
	PublicKey string `json:"publicKey"` // RSA公钥
	Challenge string `json:"challenge"` // 挑战码
}

// EncryptedLoginRequest 加密登录请求包装
type EncryptedLoginRequest struct {
	EncryptedData string `json:"encrypted_data" binding:"required"` // RSA加密后的数据
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
