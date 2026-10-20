package response

import "github.com/tanganyu1114/heimdallr-reborn/server/model/request"

// ValidateTokenResponse represents the response for token validation
type ValidateTokenResponse struct {
	Claims *request.CustomClaims `json:"claims"`
}

// RefreshTokenResponse represents the response for token refresh
type RefreshTokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

// CheckPermissionResponse represents the response for permission check
type CheckPermissionResponse struct {
	Allowed bool `json:"allowed"`
}
