package auth

import (
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

type hostMiddleware struct {
	md *authMiddleware
}

func (h *hostMiddleware) Get() httpclientv1.ClientBuilder[metav1.IDOptions, *v1.Host] {
	return applyAuthOptions(h.md, h.md.txp.Hosts().Get())
}

func (h *hostMiddleware) List() httpclientv1.ClientBuilder[metav1.ListOptions, *v1.HostList] {
	return applyAuthOptions(h.md, h.md.txp.Hosts().List())
}

func newHostMiddleware(a *authMiddleware) *hostMiddleware {
	return &hostMiddleware{md: a}
}
