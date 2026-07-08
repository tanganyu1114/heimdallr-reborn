package auth

import (
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	modelclientv1 "gin-vue-admin/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

type groupMiddleware struct {
	md *authMiddleware
}

func (g *groupMiddleware) Get() httpclientv1.ClientBuilder[metav1.IDOptions, modelclientv1.ResponseBody[*v1.Group]] {
	return wrapWithAuth(g.md, g.md.txp.Groups().Get())
}

func (g *groupMiddleware) List() httpclientv1.ClientBuilder[metav1.ListOptions, modelclientv1.ResponseBody[*v1.GroupList]] {
	return wrapWithAuth(g.md, g.md.txp.Groups().List())
}

func newGroupMiddleware(a *authMiddleware) *groupMiddleware {
	return &groupMiddleware{md: a}
}
