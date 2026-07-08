package auth

import (
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

type hostMiddleware struct {
	md *authMiddleware
}

func (h *hostMiddleware) Get() httpclientv1.ClientBuilder[metav1.IDOptions, modelclientv1.ResponseBody[*v1.Host]] {
	return wrapWithAuth(h.md, h.md.txp.Hosts().Get())
}

func (h *hostMiddleware) List() httpclientv1.ClientBuilder[metav1.ListOptions, modelclientv1.ResponseBody[*v1.HostList]] {
	return wrapWithAuth(h.md, h.md.txp.Hosts().List())
}

func newHostMiddleware(a *authMiddleware) *hostMiddleware {
	return &hostMiddleware{md: a}
}
