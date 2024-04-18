package v1

//go:generate mockgen -source=web_server_config.go -destination=mock_web_server_config.go -package=v1 gin-vue-admin/internal/hmdr_api/store/v1 Factory,AgentInfoStore,GroupStore,HostStore,WebServerConfigStore,WebServerLogWatcherStore,WebServerStatisticsStore

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type WebServerConfigStore interface {
	GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error)
	GetConfig(ctx context.Context, opts metav1.WebServerOptions) ([]string, error)
	InsertWithClone(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error
	InsertWithNew(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error
	Remove(ctx context.Context, opts metav1.WebServerOptions, pos metav1.ConfigContextPos) error
	ModifyWithClone(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error
	ModifyWithNew(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error
}
