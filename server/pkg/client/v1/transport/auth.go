package transport

import (
	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	http_transport "github.com/go-kit/kit/transport/http"
)

// AuthTransport defines the interface for authentication and authorization related transport
type AuthTransport interface {
	// ValidateToken returns the token validation client
	ValidateToken() httpclientv1.ClientBuilder[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]]
	// RefreshToken returns the token refresh client
	RefreshToken() httpclientv1.ClientBuilder[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]]
	// CheckPermission returns the permission check client
	CheckPermission() httpclientv1.ClientBuilder[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]]
	// LogOperation returns the operation log client
	LogOperation() httpclientv1.ClientBuilder[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]]
}

// authTransport implements AuthTransport interface
type authTransport struct {
	validateTokenClient   httpclientv1.ClientBuilder[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]]
	refreshTokenClient    httpclientv1.ClientBuilder[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]]
	checkPermissionClient httpclientv1.ClientBuilder[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]]
	logOperationClient    httpclientv1.ClientBuilder[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]]
}

// newAuthTransport creates a new auth transport
func newAuthTransport(transport *transport) AuthTransport {
	t := &authTransport{
		validateTokenClient: httpclientv1.NewClientBuilder[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/auth/validateToken",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		refreshTokenClient: httpclientv1.NewClientBuilder[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/auth/refreshToken",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		checkPermissionClient: httpclientv1.NewClientBuilder[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/auth/checkPermission",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		logOperationClient: httpclientv1.NewClientBuilder[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/auth/logOperation",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
	}
	return t
}

// ValidateToken returns the token validation client
func (a *authTransport) ValidateToken() httpclientv1.ClientBuilder[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]] {
	return a.validateTokenClient
}

// RefreshToken returns the token refresh client
func (a *authTransport) RefreshToken() httpclientv1.ClientBuilder[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]] {
	return a.refreshTokenClient
}

// CheckPermission returns the permission check client
func (a *authTransport) CheckPermission() httpclientv1.ClientBuilder[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]] {
	return a.checkPermissionClient
}

// LogOperation returns the operation log client
func (a *authTransport) LogOperation() httpclientv1.ClientBuilder[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	return a.logOperationClient
}
