package cache

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/marmotedu/errors"
)

type webServerConfigStore struct {
	cacheStore *cacheStore
}

func (w *webServerConfigStore) GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error) {
	return w.cacheStore.store.WebServerConfigs().GetOptions(ctx)
}

func (w *webServerConfigStore) GetConfig(ctx context.Context, opts metav1.WebServerOptions) (configuration.NginxConfig, error) {
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

func (w *webServerConfigStore) GetContext(ctx context.Context, opts metav1.WebServerOptions, pos metav1.ConfigContextPos) (nginx_context.Context, error) {
	config, err := w.cacheStore.GetConfig(ctx, opts)
	if err != nil {
		return nginx_context.NullContext(), err
	}
	nginxCtx, err := w.parseContext(config, pos.Config, pos.ContextPosPath)
	if err != nil {
		return nginxCtx, errors.Errorf("failed to parse the target context: %v", err)
	}
	return nginxCtx, nil
}

func (w *webServerConfigStore) GetIncludedConfigs(ctx context.Context, opts metav1.WebServerOptions, pos metav1.ConfigContextPos) ([]string, error) {
	c, err := w.GetContext(ctx, opts, pos)
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

func (w *webServerConfigStore) InsertWithClone(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	err := w.cacheStore.CheckBeforeUpdating(ctx, opts)
	if err != nil {
		return err
	}
	return w.cacheStore.store.WebServerConfigs().InsertWithClone(ctx, opts, ctxmeta)
}

func (w *webServerConfigStore) InsertWithNew(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	err := w.cacheStore.CheckBeforeUpdating(ctx, opts)
	if err != nil {
		return err
	}
	return w.cacheStore.store.WebServerConfigs().InsertWithNew(ctx, opts, ctxmeta)
}

func (w *webServerConfigStore) Remove(ctx context.Context, opts metav1.WebServerOptions, pos metav1.ConfigContextPos) error {
	err := w.cacheStore.CheckBeforeUpdating(ctx, opts)
	if err != nil {
		return err
	}
	return w.cacheStore.store.WebServerConfigs().Remove(ctx, opts, pos)
}

func (w *webServerConfigStore) ModifyWithClone(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	err := w.cacheStore.CheckBeforeUpdating(ctx, opts)
	if err != nil {
		return err
	}
	return w.cacheStore.store.WebServerConfigs().ModifyWithClone(ctx, opts, ctxmeta)
}

func (w *webServerConfigStore) ModifyWithNew(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	err := w.cacheStore.CheckBeforeUpdating(ctx, opts)
	if err != nil {
		return err
	}
	return w.cacheStore.store.WebServerConfigs().ModifyWithNew(ctx, opts, ctxmeta)
}

func (w *webServerConfigStore) ChangeContextEnabledState(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]) error {
	err := w.cacheStore.CheckBeforeUpdating(ctx, opts)
	if err != nil {
		return err
	}
	return w.cacheStore.store.WebServerConfigs().ChangeContextEnabledState(ctx, opts, ctxmeta)
}

func (w *webServerConfigStore) ModifyContextValue(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	err := w.cacheStore.CheckBeforeUpdating(ctx, opts)
	if err != nil {
		return err
	}
	return w.cacheStore.store.WebServerConfigs().ModifyContextValue(ctx, opts, ctxmeta)
}

func (w *webServerConfigStore) Move(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	err := w.cacheStore.CheckBeforeUpdating(ctx, opts)
	if err != nil {
		return err
	}
	return w.cacheStore.store.WebServerConfigs().Move(ctx, opts, ctxmeta)
}

func newWebServerConfigStore(cacheStore *cacheStore) *webServerConfigStore {
	return &webServerConfigStore{cacheStore: cacheStore}
}
