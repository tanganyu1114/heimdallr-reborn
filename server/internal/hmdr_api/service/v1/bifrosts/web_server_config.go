package bifrosts

import (
	"context"
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/service/v1"
	storev1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"

	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
)

type webServerConfigService struct {
	store storev1.Factory
}

var _ svcv1.WebServerConfigSrv = (*webServerConfigService)(nil)

func (w *webServerConfigService) GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error) {
	return w.store.WebServerConfigs().GetOptions(ctx)
}

func (w *webServerConfigService) GetConfig(ctx context.Context, opts metav1.WebServerOptions) (metav1.WebServerConfig, error) {
	config, ofp, err := w.store.WebServerConfigs().GetConfig(ctx, opts)
	if err != nil {
		return metav1.WebServerConfig{}, err
	}
	return metav1.WebServerConfig{
		Config:               config.Main(),
		OriginalFingerprints: ofp.Fingerprints(),
	}, nil
}

func (w *webServerConfigService) GetContext(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos metav1.ConfigContextPos) (nginx_context.Context, error) {
	return w.store.WebServerConfigs().GetContext(ctx, opts, ofp, pos)
}

func (w *webServerConfigService) GetIncludedConfigs(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos metav1.ConfigContextPos) ([]string, error) {
	return w.store.WebServerConfigs().GetIncludedConfigs(ctx, opts, ofp, pos)
}

func (w *webServerConfigService) SearchContextPositions(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, kwmeta metav1.SearchKeywordsMeta) ([]metav1.ConfigContextPos, error) {
	return w.store.WebServerConfigs().SearchContextPositions(ctx, opts, ofp, kwmeta)
}

func (w *webServerConfigService) InsertWithClone(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta], disabledTarget bool) error {
	return w.store.WebServerConfigs().InsertWithClone(ctx, opts, ofp, ctxmeta, disabledTarget)
}

func (w *webServerConfigService) InsertWithNew(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta], disabledTarget bool) error {
	return w.store.WebServerConfigs().InsertWithNew(ctx, opts, ofp, ctxmeta, disabledTarget)
}

func (w *webServerConfigService) Remove(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos metav1.ConfigContextPos) error {
	return w.store.WebServerConfigs().Remove(ctx, opts, ofp, pos)
}

func (w *webServerConfigService) UpdateConfig(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, configJsonData []byte) error {
	return w.store.WebServerConfigs().UpdateConfig(ctx, opts, ofp, configJsonData)
}

func (w *webServerConfigService) ModifyWithClone(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	return w.store.WebServerConfigs().ModifyWithClone(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigService) ModifyWithNew(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	return w.store.WebServerConfigs().ModifyWithNew(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigService) ChangeContextEnabledState(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]) error {
	return w.store.WebServerConfigs().ChangeContextEnabledState(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigService) ModifyContextValue(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	return w.store.WebServerConfigs().ModifyContextValue(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigService) Move(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta], disabledTarget bool) error {
	return w.store.WebServerConfigs().Move(ctx, opts, ofp, ctxmeta, disabledTarget)
}

func newWebServerConfigs(svc *service) svcv1.WebServerConfigSrv {
	return &webServerConfigService{
		store: svc.store,
	}
}
