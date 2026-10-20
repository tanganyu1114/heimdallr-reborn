package auth

import (
	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"
)

type authsMiddleware struct {
	md *authMiddleware
}

func (a *authsMiddleware) ValidateToken() httpclientv1.ClientBuilder[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]] {
	return wrapWithAuth(a.md, a.md.txp.Authz().ValidateToken())
}

func (a *authsMiddleware) RefreshToken() httpclientv1.ClientBuilder[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]] {
	return wrapWithAuth(a.md, a.md.txp.Authz().RefreshToken())
}

func (a *authsMiddleware) CheckPermission() httpclientv1.ClientBuilder[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]] {
	return wrapWithAuth(a.md, a.md.txp.Authz().CheckPermission())
}

func (a *authsMiddleware) LogOperation() httpclientv1.ClientBuilder[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	return wrapWithAuth(a.md, a.md.txp.Authz().LogOperation())
}

func newAuthzMiddleware(md *authMiddleware) *authsMiddleware {
	return &authsMiddleware{md: md}
}
