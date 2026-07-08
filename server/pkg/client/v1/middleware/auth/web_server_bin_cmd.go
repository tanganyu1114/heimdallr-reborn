package auth

import (
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	modelclientv1 "gin-vue-admin/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

type webServerBinCMDMiddleware struct {
	md *authMiddleware
}

func (w *webServerBinCMDMiddleware) Exec() httpclientv1.ClientBuilder[metav1.WebServerBinCMDExecRequest, modelclientv1.ResponseBody[*metav1.WebServerBinCMDExecResponse]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerBinCMDs().Exec())
}

func newWebServerBinCMDMiddleware(a *authMiddleware) *webServerBinCMDMiddleware {
	return &webServerBinCMDMiddleware{md: a}
}
