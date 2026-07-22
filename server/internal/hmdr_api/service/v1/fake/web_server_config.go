package fake

import (
	"context"

	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	storefake "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1/fake"
)

type WebServerConfigService struct {
}

func (w WebServerConfigService) GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error) {
	//TODO implement me
	panic("implement me")
}

func (w WebServerConfigService) GetConfig(ctx context.Context, opts v1.WebServerOptions) (v1.WebServerConfig, error) {
	config, ofp, err := new(storefake.WebServerConfigStore).GetConfig(ctx, opts)
	if err != nil {
		return v1.WebServerConfig{}, err
	}
	return v1.WebServerConfig{
		Config:               config.Main(),
		OriginalFingerprints: ofp.Fingerprints(),
	}, nil
}

func (w WebServerConfigService) GetContext(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos v1.ConfigContextPos) (nginx_context.Context, error) {
	return new(storefake.WebServerConfigStore).GetContext(ctx, opts, ofp, pos)
}

func (w WebServerConfigService) GetIncludedConfigs(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos v1.ConfigContextPos) ([]string, error) {
	return new(storefake.WebServerConfigStore).GetIncludedConfigs(ctx, opts, ofp, pos)
}

func (w WebServerConfigService) SearchContextPositions(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, kwmeta v1.SearchKeywordsMeta) ([]v1.ConfigContextPos, error) {
	return new(storefake.WebServerConfigStore).SearchContextPositions(ctx, opts, ofp, kwmeta)
}

func (w WebServerConfigService) InsertWithClone(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.CloneConfigContextMeta], disabledTarget bool) error {
	return new(storefake.WebServerConfigStore).InsertWithClone(ctx, opts, ofp, ctxmeta, disabledTarget)
}

func (w WebServerConfigService) InsertWithNew(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.NewConfigContextMeta], disabledTarget bool) error {
	return new(storefake.WebServerConfigStore).InsertWithNew(ctx, opts, ofp, ctxmeta, disabledTarget)
}

func (w WebServerConfigService) Remove(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos v1.ConfigContextPos) error {
	return new(storefake.WebServerConfigStore).Remove(ctx, opts, ofp, pos)
}

func (w WebServerConfigService) UpdateConfig(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, configJsonData []byte) error {
	return new(storefake.WebServerConfigStore).UpdateConfig(ctx, opts, ofp, configJsonData)
}

func (w WebServerConfigService) ModifyWithClone(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.CloneConfigContextMeta]) error {
	return new(storefake.WebServerConfigStore).ModifyWithClone(ctx, opts, ofp, ctxmeta)
}

func (w WebServerConfigService) ModifyWithNew(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.NewConfigContextMeta]) error {
	return new(storefake.WebServerConfigStore).ModifyWithNew(ctx, opts, ofp, ctxmeta)
}

func (w WebServerConfigService) ChangeContextEnabledState(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.ConfigContextEnabledStateMeta]) error {
	return new(storefake.WebServerConfigStore).ChangeContextEnabledState(ctx, opts, ofp, ctxmeta)
}

func (w WebServerConfigService) ModifyContextValue(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.NewConfigContextMeta]) error {
	return new(storefake.WebServerConfigStore).ModifyContextValue(ctx, opts, ofp, ctxmeta)
}

func (w WebServerConfigService) Move(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.CloneConfigContextMeta], disabledTarget bool) error {
	return new(storefake.WebServerConfigStore).Move(ctx, opts, ofp, ctxmeta, disabledTarget)
}
