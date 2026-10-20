package request

// ValidateTokenRequest represents the request for token validation
type ValidateTokenRequest struct {
	Token string `json:"token"`
}

// RefreshTokenRequest represents the request for token refresh
type RefreshTokenRequest struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

// CheckPermissionRequest represents the request for permission check
type CheckPermissionRequest struct {
	Sub string `json:"sub"` // subject (role/authority ID)
	Obj string `json:"obj"` // object (request URI)
	Act string `json:"act"` // action (HTTP method)
}

// LogOperationRequest represents the request for operation logging
type LogOperationRequest struct {
	IP           string `json:"ip"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	Status       int    `json:"status"`
	Latency      int64  `json:"latency"` // in nanoseconds
	Agent        string `json:"agent"`
	ErrorMessage string `json:"error_message"`
	Body         string `json:"body"`
	Resp         string `json:"resp"`
	UserID       int    `json:"user_id"`
	IsSDKLogin   bool   `json:"is_sdk_login"`
}
