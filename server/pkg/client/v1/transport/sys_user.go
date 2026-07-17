package transport

import (
	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	http_transport "github.com/go-kit/kit/transport/http"
)

// SysUserTransport defines the interface for sys user related transport
type SysUserTransport interface {
	// GetSDKChallenge returns the SDK challenge client
	GetSDKChallenge() httpclientv1.ClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]]
	// SDKLogin returns the SDK login client (accepts encrypted request)
	SDKLogin() httpclientv1.ClientBuilder[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]]
}

// sysUserTransport implements SysUserTransport interface
type sysUserTransport struct {
	sdkChallengeClient httpclientv1.ClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]]
	sdkLoginClient     httpclientv1.ClientBuilder[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]]
}

// newSysUserTransport creates a new sysUser transport
func newSysUserTransport(transport *transport) SysUserTransport {
	t := &sysUserTransport{
		sdkChallengeClient: httpclientv1.NewClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/base/sdkChallenge",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		sdkLoginClient: httpclientv1.NewClientBuilder[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/base/sdkLogin",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
	}
	return t
}

// GetSDKChallenge returns the SDK challenge client
func (s *sysUserTransport) GetSDKChallenge() httpclientv1.ClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]] {
	return s.sdkChallengeClient
}

// SDKLogin returns the SDK login client (accepts encrypted request)
func (s *sysUserTransport) SDKLogin() httpclientv1.ClientBuilder[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]] {
	return s.sdkLoginClient
}
