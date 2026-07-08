package endpoint

import (
	"github.com/tanganyu1114/heimdallr-reborn/api/heimdallr_api/v1"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
	txpclientv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/transport"
	"sync"

	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

// WebServerStatisticsEndpoints defines the interface for web server statistics related endpoints
type WebServerStatisticsEndpoints interface {
	// GetProxyServiceInfo returns the get proxy service info endpoint
	GetProxyServiceInfo() httpclientv1.Endpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]]
	// ConnectivityCheckOfProxyService returns the connectivity check of proxy service endpoint
	ConnectivityCheckOfProxyService() httpclientv1.Endpoint[metav1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]]
	// ExportProxyServiceInfoToExcel returns the export proxy service info to excel endpoint
	ExportProxyServiceInfoToExcel() httpclientv1.Endpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[[]byte]]
}

// webServerStatisticsEndpoints implements WebServerStatisticsEndpoints interface
type webServerStatisticsEndpoints struct {
	transport                               txpclientv1.WebServerStatisticsTransport
	onceGetProxyServiceInfo                 sync.Once
	getProxyServiceInfoEndpoint             httpclientv1.Endpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]]
	onceConnectivityCheckOfProxyService     sync.Once
	connectivityCheckOfProxyServiceEndpoint httpclientv1.Endpoint[metav1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]]
	onceExportProxyServiceInfoToExcel       sync.Once
	exportProxyServiceInfoToExcelEndpoint   httpclientv1.Endpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[[]byte]]
}

// GetProxyServiceInfo returns the get proxy service info endpoint
func (w *webServerStatisticsEndpoints) GetProxyServiceInfo() httpclientv1.Endpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]] {
	w.onceGetProxyServiceInfo.Do(func() {
		w.getProxyServiceInfoEndpoint = w.transport.GetProxyServiceInfo().Build().Endpoint()
	})
	return w.getProxyServiceInfoEndpoint
}

// ConnectivityCheckOfProxyService returns the connectivity check of proxy service endpoint
func (w *webServerStatisticsEndpoints) ConnectivityCheckOfProxyService() httpclientv1.Endpoint[metav1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]] {
	w.onceConnectivityCheckOfProxyService.Do(func() {
		w.connectivityCheckOfProxyServiceEndpoint = w.transport.ConnectivityCheckOfProxyService().Build().Endpoint()
	})
	return w.connectivityCheckOfProxyServiceEndpoint
}

// ExportProxyServiceInfoToExcel returns the export proxy service info to excel endpoint
func (w *webServerStatisticsEndpoints) ExportProxyServiceInfoToExcel() httpclientv1.Endpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[[]byte]] {
	w.onceExportProxyServiceInfoToExcel.Do(func() {
		w.exportProxyServiceInfoToExcelEndpoint = w.transport.ExportProxyServiceInfoToExcel().Build().Endpoint()
	})
	return w.exportProxyServiceInfoToExcelEndpoint
}

// newWebServerStatisticsEndpoints creates a new WebServerStatistics endpoints
func newWebServerStatisticsEndpoints(f *factory) WebServerStatisticsEndpoints {
	return &webServerStatisticsEndpoints{
		transport:                           f.transport.WebServerStatistics(),
		onceGetProxyServiceInfo:             sync.Once{},
		onceConnectivityCheckOfProxyService: sync.Once{},
		onceExportProxyServiceInfoToExcel:   sync.Once{},
	}
}
