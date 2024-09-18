package fake

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	storefake "gin-vue-admin/internal/hmdr_api/store/v1/fake"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
)

type WebServerConfigService struct {
}

func (w WebServerConfigService) GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error) {
	//TODO implement me
	panic("implement me")
}

func (w WebServerConfigService) GetConfig(ctx context.Context, opts metav1.WebServerOptions) (configuration.NginxConfig, error) {
	return new(storefake.WebServerConfigStore).GetConfig(ctx, opts)
}

func (w WebServerConfigService) InsertWithClone(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	return new(storefake.WebServerConfigStore).InsertWithClone(ctx, opts, ctxmeta)
}

func (w WebServerConfigService) InsertWithNew(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	return new(storefake.WebServerConfigStore).InsertWithNew(ctx, opts, ctxmeta)
}

func (w WebServerConfigService) Remove(ctx context.Context, opts metav1.WebServerOptions, pos metav1.ConfigContextPos) error {
	return new(storefake.WebServerConfigStore).Remove(ctx, opts, pos)
}

func (w WebServerConfigService) ModifyWithClone(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	return new(storefake.WebServerConfigStore).ModifyWithClone(ctx, opts, ctxmeta)
}

func (w WebServerConfigService) ModifyWithNew(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	return new(storefake.WebServerConfigStore).ModifyWithNew(ctx, opts, ctxmeta)
}

func (w WebServerConfigService) ModifyContextValue(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	return new(storefake.WebServerConfigStore).ModifyContextValue(ctx, opts, ctxmeta)
}

func (w WebServerConfigService) Move(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	return new(storefake.WebServerConfigStore).Move(ctx, opts, ctxmeta)
}
