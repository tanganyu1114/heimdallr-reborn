package transport

import (
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	modelclientv1 "gin-vue-admin/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	http_transport "github.com/go-kit/kit/transport/http"
)

// WebServerStatisticsTransport defines the interface for web server statistics related transport
type WebServerStatisticsTransport interface {
	// GetProxyServiceInfo returns the get proxy service info client
	GetProxyServiceInfo() httpclientv1.ClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]]
	// ConnectivityCheckOfProxyService returns the connectivity check of proxy service client
	ConnectivityCheckOfProxyService() httpclientv1.ClientBuilder[metav1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]]
	// ExportProxyServiceInfoToExcel returns the export proxy service info to excel client
	ExportProxyServiceInfoToExcel() httpclientv1.ClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]byte]]
}

// webServerStatisticsTransport implements WebServerStatisticsTransport interface
type webServerStatisticsTransport struct {
	getProxyServiceInfoClient             httpclientv1.ClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]]
	connectivityCheckOfProxyServiceClient httpclientv1.ClientBuilder[metav1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]]
	exportProxyServiceInfoToExcelClient   httpclientv1.ClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]byte]]
}

// newWebServerStatisticsTransport creates a new WebServerStatistics transport
func newWebServerStatisticsTransport(transport *transport) WebServerStatisticsTransport {
	t := &webServerStatisticsTransport{
		getProxyServiceInfoClient: httpclientv1.NewClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/hmdr-statistics/proxy-svc-brief",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		connectivityCheckOfProxyServiceClient: httpclientv1.NewClientBuilder[metav1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/hmdr-statistics/conn-check-of-proxy-svc",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		exportProxyServiceInfoToExcelClient: httpclientv1.NewClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]byte]](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/hmdr-statistics/export-proxy-svc-excel",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
	}
	return t
}

// GetProxyServiceInfo returns the get proxy service info client
func (w *webServerStatisticsTransport) GetProxyServiceInfo() httpclientv1.ClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]] {
	return w.getProxyServiceInfoClient
}

// ConnectivityCheckOfProxyService returns the connectivity check of proxy service client
func (w *webServerStatisticsTransport) ConnectivityCheckOfProxyService() httpclientv1.ClientBuilder[metav1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]] {
	return w.connectivityCheckOfProxyServiceClient
}

// ExportProxyServiceInfoToExcel returns the export proxy service info to excel client
func (w *webServerStatisticsTransport) ExportProxyServiceInfoToExcel() httpclientv1.ClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]byte]] {
	return w.exportProxyServiceInfoToExcelClient
}
