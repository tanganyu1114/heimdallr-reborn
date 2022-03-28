package v1

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type WebServerStatisticsSrv interface {
	GetProxyServiceInfo(ctx context.Context, opts metav1.WebServerOptions) ([]v1.ProxyServiceInfo, error)
}

type webServerStatisticsService struct {
	store storev1.Factory
}

var _ WebServerStatisticsSrv = (*webServerStatisticsService)(nil)

func (w *webServerStatisticsService) GetProxyServiceInfo(ctx context.Context, opts metav1.WebServerOptions) ([]v1.ProxyServiceInfo, error) {
	return w.store.WebServerStatistics().GetProxyServiceInfo(ctx, opts)
}

func newWebServerStatistics(svc *service) WebServerStatisticsSrv {
	return &webServerStatisticsService{
		store: svc.store,
	}
}
