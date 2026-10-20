package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/heimdallr-reborn/server/global"
	"github.com/tanganyu1114/heimdallr-reborn/server/middleware"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	"github.com/tanganyu1114/heimdallr-reborn/server/service"
	"go.uber.org/zap"
)

// BlacklistChecker is a function type for checking if a token is blacklisted
type BlacklistChecker func(jwt string) bool

// DefaultBlacklistChecker uses the service.IsBlacklist function
var DefaultBlacklistChecker BlacklistChecker = service.IsBlacklist

// TokenCreator is a function type for creating JWT tokens
type TokenCreator func(claims request.CustomClaims) (string, error)

// DefaultTokenCreator uses the JWT CreateToken method
var DefaultTokenCreator TokenCreator = nil

// TokenParser is a function type for parsing JWT tokens
type TokenParser func(token string) (*request.CustomClaims, error)

// DefaultTokenParser uses the JWT ParseToken method
var DefaultTokenParser TokenParser = nil

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(c *gin.Context) {
	var req request.ValidateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	j := middleware.NewJWT()
	claims, err := j.ParseToken(req.Token)
	if err != nil {
		global.GVA_LOG.Error("token validation failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	if DefaultBlacklistChecker(req.Token) {
		response.FailWithMessage("token is blacklisted", c)
		return
	}

	response.OkWithData(response.ValidateTokenResponse{Claims: claims}, c)
}

// RefreshToken refreshes a JWT token and returns the new token
func RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	j := middleware.NewJWT()

	// Use injected parser or default
	parseToken := DefaultTokenParser
	if parseToken == nil {
		parseToken = j.ParseToken
	}

	claims, err := parseToken(req.Token)
	if err != nil {
		global.GVA_LOG.Error("token parse failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	if DefaultBlacklistChecker(req.Token) {
		response.FailWithMessage("token is blacklisted", c)
		return
	}

	claims.ExpiresAt = req.ExpiresAt

	// Use injected creator or default
	createToken := DefaultTokenCreator
	if createToken == nil {
		createToken = j.CreateToken
	}

	newToken, err := createToken(*claims)
	if err != nil {
		global.GVA_LOG.Error("token refresh failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(response.RefreshTokenResponse{Token: newToken}, c)
}

// CheckPermission checks if a user has permission to perform an action
func CheckPermission(c *gin.Context) {
	var req request.CheckPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	e := service.Casbin()
	allowed, err := e.Enforce(req.Sub, req.Obj, req.Act)
	if err != nil {
		global.GVA_LOG.Error("permission check failed", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(response.CheckPermissionResponse{Allowed: allowed}, c)
}

// LogOperation logs an operation record
func LogOperation(c *gin.Context) {
	var req request.LogOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// Log the operation record (implementation depends on your logging service)
	global.GVA_LOG.Info("operation logged",
		zap.String("ip", req.IP),
		zap.String("method", req.Method),
		zap.String("path", req.Path),
		zap.Int("userId", req.UserID),
	)

	response.OkWithMessage("operation logged successfully", c)
}
