package v1

import (
	"context"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
)

type WebServerConfigStore interface {
	GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error)
	GetConfig(ctx context.Context, opts v1.WebServerOptions) (configuration.NginxConfig, utilsV3.ConfigFingerprinter, error)
	GetContext(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos v1.ConfigContextPos) (nginx_context.Context, error)
	GetIncludedConfigs(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos v1.ConfigContextPos) ([]string, error)
	SearchContextPositions(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, kwmeta v1.SearchKeywordsMeta) ([]v1.ConfigContextPos, error)
	InsertWithClone(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.CloneConfigContextMeta], disabledTarget bool) error
	InsertWithNew(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.NewConfigContextMeta], disabledTarget bool) error
	Remove(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos v1.ConfigContextPos) error
	UpdateConfig(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, configJsonData []byte) error
	ModifyWithClone(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.CloneConfigContextMeta]) error
	ModifyWithNew(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.NewConfigContextMeta]) error
	ChangeContextEnabledState(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.ConfigContextEnabledStateMeta]) error
	ModifyContextValue(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.NewConfigContextMeta]) error
	Move(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.CloneConfigContextMeta], disabledTarget bool) error
}
