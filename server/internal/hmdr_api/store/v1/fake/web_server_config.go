package fake

import (
	"context"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	"github.com/marmotedu/errors"
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	storev1utils "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1/utils"
	"github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/bifrosts/fake"
)

type WebServerConfigStore struct {
}

func newConfigContext(meta v1.NewConfigContextMeta) nginx_context.Context {
	return local.NewContext(meta.ContextType, meta.ContextValue)
}

func (f WebServerConfigStore) GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error) {
	//TODO implement me
	panic("implement me")
}

func (f WebServerConfigStore) GetConfig(ctx context.Context, opts v1.WebServerOptions) (configuration.NginxConfig, utilsV3.ConfigFingerprinter, error) {
	return new(fake.ServiceClient).WebServerConfig().Get(opts.ServerName)
}

func (f WebServerConfigStore) getConfigAndVerifyOFP(ctx context.Context, opts v1.WebServerOptions, fp utilsV3.ConfigFingerprints) (configuration.NginxConfig, error) {
	config, ofp, err := f.GetConfig(ctx, opts)
	if err != nil {
		return nil, err
	}
	if ofp.Diff(fp) {
		return nil, errors.Wrapf(v1.ErrInconsistentFingerprints, "the config fingerprints(%v) to be checked are different from remote server's(%v)", fp, ofp.Fingerprints())
	}
	return config, nil
}

func (f WebServerConfigStore) GetContext(ctx context.Context, opts v1.WebServerOptions, fp utilsV3.ConfigFingerprints, pos v1.ConfigContextPos) (nginx_context.Context, error) {
	config, err := f.getConfigAndVerifyOFP(ctx, opts, fp)
	if err != nil {
		return nginx_context.NullContext(), err
	}
	posConfigPath, err := nginx_context.NewRelConfigPath(config.Main().MainConfig().BaseDir(), pos.Config)
	if err != nil {
		return nginx_context.NullContext(), errors.Errorf("failed to parse nginx config path(%s), cased by: %s", pos.Config, err)
	}
	target := nginx_context.NullContext()
	target, err = config.Main().GetConfig(posConfigPath.FullPath())
	if err != nil {
		return nginx_context.NullContext(), errors.Errorf("failed to parse target context: %v", err)
	}
	for _, idx := range pos.ContextPosPath {
		target = target.Child(idx)
	}
	return target, target.Error()
}

func (f WebServerConfigStore) GetIncludedConfigs(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos v1.ConfigContextPos) ([]string, error) {
	c, err := f.GetContext(ctx, opts, ofp, pos)
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

func (f WebServerConfigStore) SearchContextPositions(ctx context.Context, opts v1.WebServerOptions, fp utilsV3.ConfigFingerprints, kwmeta v1.SearchKeywordsMeta) ([]v1.ConfigContextPos, error) {
	starting := nginx_context.NewPosSet()

	for _, ccp := range kwmeta.StartingPositionList {
		c, err := f.GetContext(ctx, opts, fp, ccp)
		if err != nil {
			return nil, err
		}
		starting.AppendWithPosSet(c.ChildrenPosSet())
	}

	return storev1utils.SearchContextPoses(starting, kwmeta.IsOnlyInCurrent, kwmeta.Keywords, kwmeta.IsRegexpRule)
}

func (f WebServerConfigStore) InsertWithClone(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.CloneConfigContextMeta], disabledTarget bool) error {
	if len(ctxmeta.TargetContext.ContextPosPath) == 0 {
		return errors.Errorf("failed to parse context posision to be cloned: %v", errors.New("nginx config context pos path is null"))
	}
	return nil
}

func (f WebServerConfigStore) InsertWithNew(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.NewConfigContextMeta], disabledTarget bool) error {
	if err := newConfigContext(ctxmeta.TargetContext).Error(); err != nil {
		return err
	}
	return nil
}

func (f WebServerConfigStore) UpdateConfig(_ context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, configJsonData []byte) error {
	return f.updateConfig(opts, ofp, func(conf configuration.NginxConfig) (configuration.NginxConfig, error) {
		newConf, err := configuration.NewNginxConfigFromJsonBytes(configJsonData)
		if err != nil {
			return conf, err
		}
		return newConf, nil
	})
}

func (f WebServerConfigStore) Remove(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos v1.ConfigContextPos) error {
	if len(pos.ContextPosPath) == 0 {
		return errors.Errorf("failed to parse target position: %v", errors.New("nginx config context pos path is null"))
	}
	return nil
}

func (f WebServerConfigStore) ModifyWithClone(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.CloneConfigContextMeta]) error {
	if len(ctxmeta.TargetContext.ContextPosPath) == 0 {
		return errors.Errorf("failed to parse context posision to be cloned: %v", errors.New("nginx config context pos path is null"))
	}
	return nil
}

func (f WebServerConfigStore) ModifyWithNew(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.NewConfigContextMeta]) error {
	if err := newConfigContext(ctxmeta.TargetContext).Error(); err != nil {
		return err
	}
	return nil
}

func (f WebServerConfigStore) updateConfig(opts v1.WebServerOptions, fp utilsV3.ConfigFingerprints, operfn func(conf configuration.NginxConfig) (configuration.NginxConfig, error)) error {
	nginxConfig, ofp, err := new(fake.ServiceClient).WebServerConfig().Get(opts.ServerName)
	if err != nil {
		return err
	}
	if ofp.Diff(fp) {
		return errors.Wrapf(v1.ErrInconsistentFingerprints, "the config fingerprints(%v) to be checked are different from remote server's(%v)", fp, ofp.Fingerprints())
	}
	updating, err := operfn(nginxConfig)
	if err != nil {
		return err
	}

	return new(fake.ServiceClient).WebServerConfig().Update(opts.ServerName, updating, fp)
}

func (f WebServerConfigStore) ChangeContextEnabledState(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.ConfigContextEnabledStateMeta]) error {
	return f.updateConfig(opts, ofp, func(conf configuration.NginxConfig) (configuration.NginxConfig, error) {
		target, err := storev1utils.ParseContext(conf, ctxmeta.Position.Config, ctxmeta.Position.ContextPosPath)
		if err != nil {
			return conf, errors.Errorf("failed to parse the target position: %v", err)
		}

		if ctxmeta.TargetContext.Enabled {
			target.Enable()
		} else {
			target.Disable()
		}
		return conf, nil
	})
}

func (f WebServerConfigStore) ModifyContextValue(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.NewConfigContextMeta]) error {
	return f.updateConfig(opts, ofp, func(conf configuration.NginxConfig) (configuration.NginxConfig, error) {
		targetPos, err := storev1utils.ParseContextPos(conf, ctxmeta.Position)
		if err != nil {
			return conf, errors.Errorf("failed to parse target position: %v", err)
		}
		if ctxmeta.TargetContext.ContextType != targetPos.Target().Type() {
			return conf, errors.Errorf("the target context type for modification is inconsistent. the target context type: %s, the modified context type: %s.", targetPos.Target().Type(), ctxmeta.TargetContext.ContextType)
		}
		err = targetPos.Target().SetValue(ctxmeta.TargetContext.ContextValue)
		if err != nil {
			return conf, errors.Errorf("failed to set value for target context: %v", err)
		}
		return conf, nil
	})
}

func (f WebServerConfigStore) Move(ctx context.Context, opts v1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta v1.TargetConfigContextOptions[v1.CloneConfigContextMeta], disabledTarget bool) error {
	return f.updateConfig(opts, ofp, func(conf configuration.NginxConfig) (configuration.NginxConfig, error) {
		targetPos, err := storev1utils.ParseContextPosModifyTO(conf, ctxmeta.Position)
		if err != nil {
			return conf, errors.Errorf("failed to parse target position: %v", err)
		}
		targetFather, targetIdx := targetPos.Position()
		toBeMovedPos, err := storev1utils.ParseContextPos(conf, ctxmeta.TargetContext.ConfigContextPos)
		if err != nil {
			return conf, errors.Errorf("failed to parse context posision to be moved: %v", err)
		}
		toBeMovedCtx := toBeMovedPos.Target()
		if err = toBeMovedCtx.Error(); err != nil {
			return conf, errors.Errorf("failed to parse context posision to be moved: %v", err)
		}
		toBeMovedFather, toBeMovedIdx := toBeMovedPos.Position()
		if err = toBeMovedFather.Remove(toBeMovedIdx).Error(); err != nil {
			return conf, errors.Errorf("failed to removed context to be moved: %v", err)
		}
		if err = targetFather.Insert(toBeMovedCtx, targetIdx).Error(); err != nil {
			return conf, errors.Errorf("failed to insert context to be moved: %v", err)
		}
		if disabledTarget {
			err = toBeMovedCtx.Disable().Error()
		}
		return conf, err
	})
}
