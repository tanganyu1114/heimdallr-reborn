package auth

import (
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

type webServerStatisticsMiddleware struct {
	md *authMiddleware
}

func (w *webServerStatisticsMiddleware) GetProxyServiceInfo() httpclientv1.ClientBuilder[metav1.WebServerOptions, []v1.ProxyServiceInfo] {
	return applyAuthOptions(w.md, w.md.txp.WebServerStatistics().GetProxyServiceInfo())
}

func (w *webServerStatisticsMiddleware) ConnectivityCheckOfProxyService() httpclientv1.ClientBuilder[metav1.ConnectivityCheckOfProxiedServersRequestOptions, v1.ProxyServiceInfo] {
	return applyAuthOptions(w.md, w.md.txp.WebServerStatistics().ConnectivityCheckOfProxyService())
}

func (w *webServerStatisticsMiddleware) ExportProxyServiceInfoToExcel() httpclientv1.ClientBuilder[metav1.WebServerOptions, []byte] {
	return applyAuthOptions(w.md, w.md.txp.WebServerStatistics().ExportProxyServiceInfoToExcel())
}

func newWebServerStatisticsMiddleware(a *authMiddleware) *webServerStatisticsMiddleware {
	return &webServerStatisticsMiddleware{md: a}
}
