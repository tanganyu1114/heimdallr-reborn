package v1

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type WebServerConfigSrv interface {
	GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error)
	GetConfig(ctx context.Context, opts metav1.WebServerOptions) ([]byte, error)
}

type webServerConfigService struct {
	store storev1.Factory
}

var _ WebServerConfigSrv = (*webServerConfigService)(nil)

func (w *webServerConfigService) GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error) {
	return w.store.WebServerConfigs().GetOptions(ctx)
}

func (w *webServerConfigService) GetConfig(ctx context.Context, opts metav1.WebServerOptions) ([]byte, error) {
	return w.store.WebServerConfigs().GetConfig(ctx, opts)
}

func newWebServerConfigs(svc *service) WebServerConfigSrv {
	return &webServerConfigService{
		store: svc.store,
	}
}
