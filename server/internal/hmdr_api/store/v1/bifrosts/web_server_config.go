package bifrosts

import (
	"context"
	v1 "github.com/tanganyu1114/heimdallr-reborn/api/heimdallr_api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/global"
	storev1utils "github.com/tanganyu1114/heimdallr-reborn/internal/hmdr_api/store/v1/utils"
	"github.com/tanganyu1114/heimdallr-reborn/internal/pkg/bifrosts"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
	"github.com/tanganyu1114/heimdallr-reborn/pkg/sort_map"
	"sync"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	"github.com/marmotedu/errors"
	"go.uber.org/zap"
)

var (
	wscOnce           = new(sync.Once)
	singletonWSCStore *webServerConfigStore
)

type webServerConfigStore struct {
	bm bifrosts.Manager
}

func (w *webServerConfigStore) GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error) {
	var result []v1.BifrostGroupMeta
	getOptsFromGroupFunc := func(keyer sort_map.Keyer, v interface{}) bool {
		group := v.(*bifrosts.Group)
		tmpGroup := v1.BifrostGroupMeta{
			UintObjectMeta: metav1.UintObjectMeta{
				Label: group.Meta.Name,
				Value: group.Meta.ID,
			},
			Children: make([]*v1.BifrostMeta, 0),
		}

		getOptFromBifrostFunc := func(keyer sort_map.Keyer, v interface{}) bool {
			bifrost := v.(*bifrosts.Bifrost)
			tmpHost := v1.BifrostMeta{
				UintObjectMeta: metav1.UintObjectMeta{
					Label: bifrost.Meta.Name,
					Value: bifrost.Meta.ID,
				},
				Children: make([]*v1.WebSrvMeta, 0),
			}

			metrics, err := bifrost.Client.WebServerStatus().Get()
			if err != nil {
				global.GVA_LOG.Error("failed to get the web server status", zap.String("err", err.Error()))
				//return false  // 注意：如果报错，散列表后续元素将不被加载
				return true
			}
			for _, srvInfo := range metrics.StatusList {
				tmpHost.Children = append(tmpHost.Children,
					&v1.WebSrvMeta{
						StringObjectMeta: metav1.StringObjectMeta{
							Label: srvInfo.Name,
							Value: srvInfo.Name,
						},
					},
				)
			}
			tmpGroup.Children = append(tmpGroup.Children, &tmpHost)
			return true
		}

		group.Bifrosts.Range(getOptFromBifrostFunc)
		result = append(result, tmpGroup)
		return true
	}

	w.bm.Range(getOptsFromGroupFunc)
	return result, nil
}

func (w *webServerConfigStore) getConfig(opts metav1.WebServerOptions) (configuration.NginxConfig, utilsV3.ConfigFingerprinter, error) {
	client, err := w.bm.GetBifrostClient(opts)
	if err != nil {
		return nil, nil, err
	}
	nginxConfig, ofp, err := client.WebServerConfig().Get(opts.ServerName)
	if err != nil {
		return nginxConfig, ofp, errors.Wrapf(err, "failed to parse the web server config(%s)", opts.ServerName)
	}
	return nginxConfig, ofp, nil
}

func (w *webServerConfigStore) GetConfig(_ context.Context, opts metav1.WebServerOptions) (configuration.NginxConfig, utilsV3.ConfigFingerprinter, error) {
	return w.getConfig(opts)
}

func (w *webServerConfigStore) getConfigAndVerifyOFP(ctx context.Context, opts metav1.WebServerOptions, fp utilsV3.ConfigFingerprints) (configuration.NginxConfig, error) {
	config, ofp, err := w.getConfig(opts)
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
	nginxCtx, err := storev1utils.ParseContext(config, pos.Config, pos.ContextPosPath)
	if err != nil {
		return nginxCtx, errors.Errorf("failed to parse the target context: %v", err)
	}
	return nginxCtx, nil
}

func (w *webServerConfigStore) GetIncludedConfigs(ctx context.Context, opts metav1.WebServerOptions, fp utilsV3.ConfigFingerprints, pos metav1.ConfigContextPos) ([]string, error) {
	c, err := w.GetContext(ctx, opts, fp, pos)
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

func (w *webServerConfigStore) updateConfig(opts metav1.WebServerOptions, fp utilsV3.ConfigFingerprints, operfn func(conf configuration.NginxConfig) (configuration.NginxConfig, error)) error {
	nginxConfig, ofp, err := w.getConfig(opts)
	if err != nil {
		return err
	}
	if ofp.Diff(fp) {
		global.GVA_LOG.Warn("the config fingerprints to be checked are different from remote server's", zap.Any("checking", fp), zap.Any("remote", ofp.Fingerprints()))
		return errors.Wrapf(metav1.ErrInconsistentFingerprints, "the config fingerprints(%v) to be checked are different from remote server's(%v)", fp, ofp.Fingerprints())
	}
	updating, err := operfn(nginxConfig)
	if err != nil {
		return err
	}
	client, err := w.bm.GetBifrostClient(opts)
	if err != nil {
		return err
	}
	return client.WebServerConfig().Update(opts.ServerName, updating, fp)
}

func (w *webServerConfigStore) InsertWithClone(_ context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta], disabledTarget bool) error {
	return w.updateConfig(opts, ofp, func(conf configuration.NginxConfig) (configuration.NginxConfig, error) {
		targetPos, err := storev1utils.ParseContextPosModifyTO(conf, ctxmeta.Position)
		if err != nil {
			return conf, errors.Errorf("failed to parse the target position: %v", err)
		}
		father, targetIdx := targetPos.Position()

		toBeClonedPos, err := storev1utils.ParseContextPos(conf, ctxmeta.TargetContext.ConfigContextPos)
		if err != nil {
			return conf, errors.Errorf("failed to parse the context posision to be cloned: %v", err)
		}
		targetCtx := father.Insert(toBeClonedPos.Target().Clone(), targetIdx).Child(targetIdx)
		if disabledTarget {
			targetCtx.Disable()
		}
		return conf, targetCtx.Error()
	})
}

func (w *webServerConfigStore) InsertWithNew(_ context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta], disabledTarget bool) error {
	return w.updateConfig(opts, ofp, func(conf configuration.NginxConfig) (configuration.NginxConfig, error) {
		targetPos, err := storev1utils.ParseContextPosModifyTO(conf, ctxmeta.Position)
		if err != nil {
			return conf, errors.Errorf("failed to parse the target position: %v", err)
		}
		father, targetIdx := targetPos.Position()
		targetCtx := father.Insert(newConfigContext(ctxmeta.TargetContext), targetIdx).Child(targetIdx)
		if disabledTarget {
			targetCtx.Disable()
		}
		return conf, targetCtx.Error()
	})
}

func (w *webServerConfigStore) Remove(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos metav1.ConfigContextPos) error {
	return w.updateConfig(opts, ofp, func(conf configuration.NginxConfig) (configuration.NginxConfig, error) {
		targetPos, err := storev1utils.ParseContextPos(conf, pos)
		if err != nil {
			return conf, errors.Errorf("failed to parse the target position: %v", err)
		}
		father, targetIdx := targetPos.Position()
		return conf, father.Remove(targetIdx).Error()
	})
}

func (w *webServerConfigStore) UpdateConfig(_ context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, configJsonData []byte) error {
	return w.updateConfig(opts, ofp, func(conf configuration.NginxConfig) (configuration.NginxConfig, error) {
		newConf, err := configuration.NewNginxConfigFromJsonBytes(configJsonData)
		if err != nil {
			return conf, err
		}
		return newConf, nil
	})
}
func (w *webServerConfigStore) ModifyWithClone(_ context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	return w.updateConfig(opts, ofp, func(conf configuration.NginxConfig) (configuration.NginxConfig, error) {
		targetPos, err := storev1utils.ParseContextPos(conf, ctxmeta.Position)
		if err != nil {
			return conf, errors.Errorf("failed to parse the target position: %v", err)
		}
		father, targetIdx := targetPos.Position()

		toBeClonedPos, err := storev1utils.ParseContextPos(conf, ctxmeta.TargetContext.ConfigContextPos)
		if err != nil {
			return conf, errors.Errorf("failed to parse the context posision to be cloned: %v", err)
		}

		return conf, father.Modify(toBeClonedPos.Target().Clone(), targetIdx).Error()
	})
}

func (w *webServerConfigStore) ModifyWithNew(_ context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	return w.updateConfig(opts, ofp, func(conf configuration.NginxConfig) (configuration.NginxConfig, error) {
		targetPos, err := storev1utils.ParseContextPos(conf, ctxmeta.Position)
		if err != nil {
			return conf, errors.Errorf("failed to parse the target position: %v", err)
		}
		father, targetIdx := targetPos.Position()
		return conf, father.Modify(newConfigContext(ctxmeta.TargetContext), targetIdx).Error()
	})
}

func (w *webServerConfigStore) ChangeContextEnabledState(_ context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]) error {
	return w.updateConfig(opts, ofp, func(conf configuration.NginxConfig) (configuration.NginxConfig, error) {
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

func (w *webServerConfigStore) ModifyContextValue(_ context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	return w.updateConfig(opts, ofp, func(conf configuration.NginxConfig) (configuration.NginxConfig, error) {
		targetPos, err := storev1utils.ParseContextPos(conf, ctxmeta.Position)
		if err != nil {
			return conf, errors.Errorf("failed to parse the target position: %v", err)
		}
		if ctxmeta.TargetContext.ContextType != targetPos.Target().Type() {
			return conf, errors.Errorf("the target context type for modification is inconsistent. the target context type: %s, the modified context type: %s.", targetPos.Target().Type(), ctxmeta.TargetContext.ContextType)
		}
		err = targetPos.Target().SetValue(ctxmeta.TargetContext.ContextValue)
		if err != nil {
			return conf, errors.Errorf("failed to set value for the target context: %v", err)
		}
		return conf, nil
	})
}

func (w *webServerConfigStore) Move(_ context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta], disabledTarget bool) error {
	return w.updateConfig(opts, ofp, func(conf configuration.NginxConfig) (configuration.NginxConfig, error) {
		targetPos, err := storev1utils.ParseContextPosModifyTO(conf, ctxmeta.Position)
		if err != nil {
			return conf, errors.Errorf("failed to parse the target position: %v", err)
		}
		targetFather, targetIdx := targetPos.Position()
		toBeMovedPos, err := storev1utils.ParseContextPos(conf, ctxmeta.TargetContext.ConfigContextPos)
		if err != nil {
			return conf, errors.Errorf("failed to parse the context posision to be moved: %v", err)
		}
		toBeMovedCtx := toBeMovedPos.Target()
		if err = toBeMovedCtx.Error(); err != nil {
			return conf, errors.Errorf("failed to parse the context posision to be moved: %v", err)
		}
		toBeMovedFather, toBeMovedIdx := toBeMovedPos.Position()
		maybeAfterTarget := toBeMovedIdx < targetIdx
		if !maybeAfterTarget {
			if err = toBeMovedFather.Remove(toBeMovedIdx).Error(); err != nil {
				return conf, errors.Errorf("failed to removed the context to be moved: %v", err)
			}
		}
		if err = targetFather.Insert(toBeMovedCtx, targetIdx).Error(); err != nil {
			return conf, errors.Errorf("failed to insert the context to be moved: %v", err)
		}
		if maybeAfterTarget {
			if err = toBeMovedFather.Remove(toBeMovedIdx).Error(); err != nil {
				return conf, errors.Errorf("failed to removed the context to be moved: %v", err)
			}
		}
		if disabledTarget {
			err = toBeMovedCtx.Disable().Error()
		}
		return conf, err
	})
}

func newWebServerConfigStore(store *bifrostsStore) *webServerConfigStore {
	wscOnce.Do(func() {
		if singletonWSCStore == nil {
			singletonWSCStore = &webServerConfigStore{
				bm: store.bm,
			}
		}
	})
	if singletonWSCStore == nil {
		panic(errors.New("the singleton web server config store is nil"))
	}
	return singletonWSCStore
}

func newConfigContext(meta metav1.NewConfigContextMeta) nginx_context.Context {
	ctx := local.NewContext(meta.ContextType, meta.ContextValue)
	if !meta.Enabled {
		ctx.Disable()
	}
	for idx, childMeta := range meta.ChildrenContextMeta {
		ctx = ctx.Insert(newConfigContext(childMeta), idx)
	}
	return ctx
}
