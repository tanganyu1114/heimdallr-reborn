package transport

import (
	"gin-vue-admin/model/request"
	"gin-vue-admin/model/response"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	http_transport "github.com/go-kit/kit/transport/http"
)

// SysUserTransport defines the interface for sys user related transport
type SysUserTransport interface {
	// SDKLogin returns the SDK login client
	SDKLogin() httpclientv1.ClientBuilder[*request.SDKLogin, *response.LoginResponse]
}

// sysUserTransport implements SysUserTransport interface
type sysUserTransport struct {
	sdkLoginClient httpclientv1.ClientBuilder[*request.SDKLogin, *response.LoginResponse]
}

// newSysUserTransport creates a new sysUser transport
func newSysUserTransport(transport *transport) SysUserTransport {
	t := &sysUserTransport{
		sdkLoginClient: httpclientv1.NewClientBuilder[*request.SDKLogin, *response.LoginResponse](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/base/sdkLogin",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
	}
	return t
}

// SDKLogin returns the SDK login client
func (s *sysUserTransport) SDKLogin() httpclientv1.ClientBuilder[*request.SDKLogin, *response.LoginResponse] {
	return s.sdkLoginClient
}
