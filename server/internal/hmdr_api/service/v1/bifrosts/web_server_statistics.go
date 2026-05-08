package bifrosts

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type webServerStatisticsService struct {
	store storev1.Factory
}

var _ svcv1.WebServerStatisticsSrv = (*webServerStatisticsService)(nil)

func (w *webServerStatisticsService) GetProxyServiceInfo(ctx context.Context, opts metav1.WebServerOptions) ([]v1.ProxyServiceInfo, error) {
	return w.store.WebServerStatistics().GetProxyServiceInfo(ctx, opts)
}

func (w *webServerStatisticsService) ConnectivityCheckOfProxyService(ctx context.Context, opts metav1.WebServerOptions, proxyPassPos metav1.ConfigContextPos) (v1.ProxyServiceInfo, error) {
	return w.store.WebServerStatistics().ConnectivityCheckOfProxyService(ctx, opts, proxyPassPos)
}

func (w *webServerStatisticsService) ExportProxyServiceInfoToExcel(ctx context.Context, opts metav1.WebServerOptions) ([]byte, error) {
	return w.store.WebServerStatistics().ExportProxyServiceInfoToExcel(ctx, opts)
}

func newWebServerStatistics(svc *service) svcv1.WebServerStatisticsSrv {
	return &webServerStatisticsService{
		store: svc.store,
	}
}
