package auth

import (
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	modelclientv1 "gin-vue-admin/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

type agentInfoMiddleware struct {
	md *authMiddleware
}

func (a *agentInfoMiddleware) Get() httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.GroupInfo]] {
	return wrapWithAuth(a.md, a.md.txp.AgentInfos().Get())
}

func newAgentInfoMiddleware(a *authMiddleware) *agentInfoMiddleware {
	return &agentInfoMiddleware{md: a}
}
