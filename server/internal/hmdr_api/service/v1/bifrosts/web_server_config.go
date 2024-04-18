package bifrosts

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type webServerConfigService struct {
	store storev1.Factory
}

var _ svcv1.WebServerConfigSrv = (*webServerConfigService)(nil)

func (w *webServerConfigService) GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error) {
	return w.store.WebServerConfigs().GetOptions(ctx)
}

func (w *webServerConfigService) GetConfig(ctx context.Context, opts metav1.WebServerOptions) ([]string, error) {
	return w.store.WebServerConfigs().GetConfig(ctx, opts)
}

func (w *webServerConfigService) InsertWithClone(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	return w.store.WebServerConfigs().InsertWithClone(ctx, opts, ctxmeta)
}

func (w *webServerConfigService) InsertWithNew(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	return w.store.WebServerConfigs().InsertWithNew(ctx, opts, ctxmeta)
}

func (w *webServerConfigService) Remove(ctx context.Context, opts metav1.WebServerOptions, pos metav1.ConfigContextPos) error {
	return w.store.WebServerConfigs().Remove(ctx, opts, pos)
}

func (w *webServerConfigService) ModifyWithClone(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	return w.store.WebServerConfigs().ModifyWithClone(ctx, opts, ctxmeta)
}

func (w *webServerConfigService) ModifyWithNew(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	return w.store.WebServerConfigs().ModifyWithNew(ctx, opts, ctxmeta)
}

func newWebServerConfigs(svc *service) svcv1.WebServerConfigSrv {
	return &webServerConfigService{
		store: svc.store,
	}
}
