package auth

import (
	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"
)

type apiMiddleware struct {
	md *authMiddleware
}

func (a *apiMiddleware) CreateApi() httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	return wrapWithAuth(a.md, a.md.txp.Apis().CreateApi())
}

func (a *apiMiddleware) UpdateApi() httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	return wrapWithAuth(a.md, a.md.txp.Apis().UpdateApi())
}

func (a *apiMiddleware) GetAllApis() httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]] {
	return wrapWithAuth(a.md, a.md.txp.Apis().GetAllApis())
}

func newApisMiddleware(md *authMiddleware) *apiMiddleware {
	return &apiMiddleware{md: md}
}
