package auth

import (
	v1 "gin-vue-admin/api/heimdallr_api/v1"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

type agentInfoMiddleware struct {
	md *authMiddleware
}

func (a *agentInfoMiddleware) Get() httpclientv1.ClientBuilder[httpclientv1.NilBody, []v1.GroupInfo] {
	return applyAuthOptions(a.md, a.md.txp.AgentInfos().Get())
}

func newAgentInfoMiddleware(a *authMiddleware) *agentInfoMiddleware {
	return &agentInfoMiddleware{md: a}
}
