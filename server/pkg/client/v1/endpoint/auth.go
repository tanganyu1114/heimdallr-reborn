package endpoint

import (
	"sync"

	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	txpclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/transport"

	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

// AuthEndpoints defines the interface for auth related endpoints
type AuthEndpoints interface {
	// ValidateToken returns the token validation endpoint
	ValidateToken() httpclientv1.Endpoint[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]]
	// RefreshToken returns the token refresh endpoint
	RefreshToken() httpclientv1.Endpoint[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]]
	// CheckPermission returns the permission check endpoint
	CheckPermission() httpclientv1.Endpoint[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]]
	// LogOperation returns the operation log endpoint
	LogOperation() httpclientv1.Endpoint[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]]
}

// authEndpoints implements AuthEndpoints interface
type authEndpoints struct {
	transport               txpclientv1.AuthTransport
	onceValidateToken       sync.Once
	validateTokenEndpoint   httpclientv1.Endpoint[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]]
	onceRefreshToken        sync.Once
	refreshTokenEndpoint    httpclientv1.Endpoint[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]]
	onceCheckPermission     sync.Once
	checkPermissionEndpoint httpclientv1.Endpoint[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]]
	onceLogOperation        sync.Once
	logOperationEndpoint    httpclientv1.Endpoint[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]]
}

// ValidateToken returns the token validation endpoint
func (a *authEndpoints) ValidateToken() httpclientv1.Endpoint[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]] {
	a.onceValidateToken.Do(func() {
		a.validateTokenEndpoint = a.transport.ValidateToken().Build().Endpoint()
	})
	return a.validateTokenEndpoint
}

// RefreshToken returns the token refresh endpoint
func (a *authEndpoints) RefreshToken() httpclientv1.Endpoint[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]] {
	a.onceRefreshToken.Do(func() {
		a.refreshTokenEndpoint = a.transport.RefreshToken().Build().Endpoint()
	})
	return a.refreshTokenEndpoint
}

// CheckPermission returns the permission check endpoint
func (a *authEndpoints) CheckPermission() httpclientv1.Endpoint[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]] {
	a.onceCheckPermission.Do(func() {
		a.checkPermissionEndpoint = a.transport.CheckPermission().Build().Endpoint()
	})
	return a.checkPermissionEndpoint
}

// LogOperation returns the operation log endpoint
func (a *authEndpoints) LogOperation() httpclientv1.Endpoint[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	a.onceLogOperation.Do(func() {
		a.logOperationEndpoint = a.transport.LogOperation().Build().Endpoint()
	})
	return a.logOperationEndpoint
}

// newAuthEndpoints creates a new Auth endpoints
func newAuthEndpoints(f *factory) AuthEndpoints {
	return &authEndpoints{
		transport: f.transport.Authz(),
	}
}
