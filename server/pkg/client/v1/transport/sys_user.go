package transport

import (
	"github.com/tanganyu1114/heimdallr-reborn/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/model/response"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	http_transport "github.com/go-kit/kit/transport/http"
)

// SysUserTransport defines the interface for sys user related transport
type SysUserTransport interface {
	// SDKLogin returns the SDK login client
	SDKLogin() httpclientv1.ClientBuilder[*request.SDKLogin, modelclientv1.ResponseBody[*response.LoginResponse]]
}

// sysUserTransport implements SysUserTransport interface
type sysUserTransport struct {
	sdkLoginClient httpclientv1.ClientBuilder[*request.SDKLogin, modelclientv1.ResponseBody[*response.LoginResponse]]
}

// newSysUserTransport creates a new sysUser transport
func newSysUserTransport(transport *transport) SysUserTransport {
	t := &sysUserTransport{
		sdkLoginClient: httpclientv1.NewClientBuilder[*request.SDKLogin, modelclientv1.ResponseBody[*response.LoginResponse]](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/base/sdkLogin",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
	}
	return t
}

// SDKLogin returns the SDK login client
func (s *sysUserTransport) SDKLogin() httpclientv1.ClientBuilder[*request.SDKLogin, modelclientv1.ResponseBody[*response.LoginResponse]] {
	return s.sdkLoginClient
}
