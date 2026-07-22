package bifrosts

import (
	"context"

	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/service/v1"
	storev1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1"
)

type webServerStatisticsService struct {
	store storev1.Factory
}

var _ svcv1.WebServerStatisticsSrv = (*webServerStatisticsService)(nil)

func (w *webServerStatisticsService) GetProxyServiceInfo(ctx context.Context, opts v1.WebServerOptions) ([]v1.ProxyServiceInfo, error) {
	return w.store.WebServerStatistics().GetProxyServiceInfo(ctx, opts)
}

func (w *webServerStatisticsService) ConnectivityCheckOfProxyService(ctx context.Context, opts v1.WebServerOptions, proxyPassPos v1.ConfigContextPos) (v1.ProxyServiceInfo, error) {
	return w.store.WebServerStatistics().ConnectivityCheckOfProxyService(ctx, opts, proxyPassPos)
}

func (w *webServerStatisticsService) ExportProxyServiceInfoToExcel(ctx context.Context, opts v1.WebServerOptions) ([]byte, error) {
	return w.store.WebServerStatistics().ExportProxyServiceInfoToExcel(ctx, opts)
}

func newWebServerStatistics(svc *service) svcv1.WebServerStatisticsSrv {
	return &webServerStatisticsService{
		store: svc.store,
	}
}
