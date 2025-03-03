package cache

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/global"
	storev1utils "gin-vue-admin/internal/hmdr_api/store/v1/utils"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	"github.com/marmotedu/errors"
	"go.uber.org/zap"
)

type webServerConfigStore struct {
	cacheStore *cacheStore
}

func (w *webServerConfigStore) GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error) {
	return w.cacheStore.next.WebServerConfigs().GetOptions(ctx)
}

func (w *webServerConfigStore) GetConfig(ctx context.Context, opts metav1.WebServerOptions) (configuration.NginxConfig, utilsV3.ConfigFingerprinter, error) {
	return w.cacheStore.GetConfig(ctx, opts)
}

func (w *webServerConfigStore) parseContext(nginxconfig configuration.NginxConfig, configPath string, ctxPosPath []int) (nginx_context.Context, error) {
	posConfigPath, err := nginx_context.NewRelConfigPath(nginxconfig.Main().MainConfig().BaseDir(), configPath)
	if err != nil {
		return nginx_context.NullContext(), errors.Errorf("failed to parse the nginx config path(%s), cased by: %s", configPath, err)
	}
	target := nginx_context.NullContext()
	target, err = nginxconfig.Main().GetConfig(posConfigPath.FullPath())
	if err != nil {
		return nginx_context.NullContext(), err
	}
	for _, idx := range ctxPosPath {
		target = target.Child(idx)
	}
	return target, target.Error()
}

func (w *webServerConfigStore) getConfigAndVerifyOFP(ctx context.Context, opts metav1.WebServerOptions, fp utilsV3.ConfigFingerprints) (configuration.NginxConfig, error) {
	config, ofp, err := w.cacheStore.GetConfig(ctx, opts)
	if err != nil {
		return nil, err
	}
	if ofp.Diff(fp) {
		global.GVA_LOG.Info("the config fingerprints to be checked are different from remote server's", zap.Any("checking", fp), zap.Any("remote", ofp.Fingerprints()))
		return nil, errors.Wrapf(metav1.ErrInconsistentFingerprints, "the config fingerprints(%v) to be checked are different from remote server's(%v)", fp, ofp.Fingerprints())
	}
	return config, nil
}

func (w *webServerConfigStore) GetContext(ctx context.Context, opts metav1.WebServerOptions, fp utilsV3.ConfigFingerprints, pos metav1.ConfigContextPos) (nginx_context.Context, error) {
	config, err := w.getConfigAndVerifyOFP(ctx, opts, fp)
	if err != nil {
		return nginx_context.NullContext(), err
	}
	nginxCtx, err := w.parseContext(config, pos.Config, pos.ContextPosPath)
	if err != nil {
		return nginxCtx, errors.Errorf("failed to parse the target context: %v", err)
	}
	return nginxCtx, nil
}

func (w *webServerConfigStore) GetIncludedConfigs(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos metav1.ConfigContextPos) ([]string, error) {
	c, err := w.GetContext(ctx, opts, ofp, pos)
	if err != nil {
		return nil, err
	}
	include, ok := c.(*local.Include)
	if !ok {
		return nil, errors.Errorf("failed to parse the target include context, possibly due to changes in the content of the target nginx config!")
	}
	includedConfigs := make([]string, 0)
	for _, config := range include.Configs() {
		includedConfigs = append(includedConfigs, config.FullPath())
	}
	return includedConfigs, nil
}

func (w *webServerConfigStore) SearchContextPositions(ctx context.Context, opts metav1.WebServerOptions, fp utilsV3.ConfigFingerprints, kwmeta metav1.SearchKeywordsMeta) ([]metav1.ConfigContextPos, error) {
	starting := nginx_context.NewPosSet()

	for _, ccp := range kwmeta.StartingPositionList {
		c, err := w.GetContext(ctx, opts, fp, ccp)
		if err != nil {
			return nil, err
		}
		starting.AppendWithPosSet(c.ChildrenPosSet())
	}

	return storev1utils.SearchContextPoses(starting, kwmeta.IsOnlyInCurrent, kwmeta.Keywords, kwmeta.IsRegexpRule)
}

func (w *webServerConfigStore) InsertWithClone(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	defer w.cacheStore.ReleaseConfigCache(opts)
	return w.cacheStore.next.WebServerConfigs().InsertWithClone(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigStore) InsertWithNew(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	defer w.cacheStore.ReleaseConfigCache(opts)
	return w.cacheStore.next.WebServerConfigs().InsertWithNew(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigStore) Remove(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos metav1.ConfigContextPos) error {
	defer w.cacheStore.ReleaseConfigCache(opts)
	return w.cacheStore.next.WebServerConfigs().Remove(ctx, opts, ofp, pos)
}

func (w *webServerConfigStore) ModifyWithClone(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	defer w.cacheStore.ReleaseConfigCache(opts)
	return w.cacheStore.next.WebServerConfigs().ModifyWithClone(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigStore) ModifyWithNew(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	defer w.cacheStore.ReleaseConfigCache(opts)
	return w.cacheStore.next.WebServerConfigs().ModifyWithNew(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigStore) ChangeContextEnabledState(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]) error {
	defer w.cacheStore.ReleaseConfigCache(opts)
	return w.cacheStore.next.WebServerConfigs().ChangeContextEnabledState(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigStore) ModifyContextValue(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	defer w.cacheStore.ReleaseConfigCache(opts)
	return w.cacheStore.next.WebServerConfigs().ModifyContextValue(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigStore) Move(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	defer w.cacheStore.ReleaseConfigCache(opts)
	return w.cacheStore.next.WebServerConfigs().Move(ctx, opts, ofp, ctxmeta)
}

func newWebServerConfigStore(cacheStore *cacheStore) *webServerConfigStore {
	return &webServerConfigStore{cacheStore: cacheStore}
}
