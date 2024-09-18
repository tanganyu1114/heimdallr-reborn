package bifrosts

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/global"
	"gin-vue-admin/internal/pkg/bifrosts"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/pkg/sort_map"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/marmotedu/errors"
	"go.uber.org/zap"
	"sync"
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
				global.GVA_LOG.Error("failed to get web server status", zap.String("err", err.Error()))
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

func (w *webServerConfigStore) getConfig(opts metav1.WebServerOptions) (configuration.NginxConfig, error) {
	data, err := w.getConfigJSON(opts)
	if err != nil {
		return nil, err
	}
	nginxConfig, err := configuration.NewNginxConfigFromJsonBytes(data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse web server config(%s)", opts.ServerName)
	}
	return nginxConfig, nil
}

func (w *webServerConfigStore) getConfigJSON(opts metav1.WebServerOptions) ([]byte, error) {
	client, err := w.bm.GetBifrostClient(opts)
	if err != nil {
		return nil, err
	}
	data, err := client.WebServerConfig().Get(opts.ServerName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get web server config(%s)", opts.ServerName)
	}
	return data, nil
}

func (w *webServerConfigStore) updateConfig(opts metav1.WebServerOptions, data []byte) error {
	if len(data) == 0 {
		return errors.New("update config data is null")
	}
	client, err := w.bm.GetBifrostClient(opts)
	if err != nil {
		return err
	}
	return client.WebServerConfig().Update(opts.ServerName, data)
}

func (w *webServerConfigStore) GetConfig(ctx context.Context, opts metav1.WebServerOptions) (configuration.NginxConfig, error) {
	return w.getConfig(opts)
}

func (w *webServerConfigStore) parseContextPos(nginxconfig configuration.NginxConfig, pos metav1.ConfigContextPos) (nginx_context.Pos, error) {
	if len(pos.ContextPosPath) == 0 {
		return nginx_context.NullPos(), errors.New("nginx config context pos path is null")
	}

	posConfigPath, err := nginx_context.NewRelConfigPath(nginxconfig.Main().MainConfig().BaseDir(), pos.Config)
	if err != nil {
		return nginx_context.NullPos(), errors.Errorf("failed to parse nginx config path(%s), cased by: %s", pos.Config, err)
	}
	nextFather := nginx_context.NullContext()
	nextFather, err = nginxconfig.Main().GetConfig(posConfigPath.FullPath())
	if err != nil {
		return nginx_context.NullPos(), err
	}
	targetIdx := pos.ContextPosPath[len(pos.ContextPosPath)-1]
	posPath := pos.ContextPosPath[:len(pos.ContextPosPath)-1]
	for _, idx := range posPath {
		nextFather = nextFather.Child(idx)
	}
	return nginx_context.SetPos(nextFather, targetIdx), nil
}

func (w *webServerConfigStore) InsertWithClone(_ context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	nginxConfig, err := w.getConfig(opts)
	if err != nil {
		return err
	}
	targetPos, err := w.parseContextPos(nginxConfig, ctxmeta.Position)
	if err != nil {
		return errors.Errorf("failed to parse target position: %v", err)
	}
	father, targetIdx := targetPos.Position()

	toBeClonedPos, err := w.parseContextPos(nginxConfig, ctxmeta.TargetContext.ConfigContextPos)
	if err != nil {
		return errors.Errorf("failed to parse context posision to be cloned: %v", err)
	}

	if err = father.Insert(toBeClonedPos.Target().Clone(), targetIdx).Error(); err != nil {
		return err
	}
	return w.updateConfig(opts, nginxConfig.Json())
}

func (w *webServerConfigStore) InsertWithNew(_ context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	nginxConfig, err := w.getConfig(opts)
	if err != nil {
		return err
	}
	targetPos, err := w.parseContextPos(nginxConfig, ctxmeta.Position)
	if err != nil {
		return errors.Errorf("failed to parse target position: %v", err)
	}
	father, targetIdx := targetPos.Position()
	if err = father.Insert(newConfigContext(ctxmeta.TargetContext), targetIdx).Error(); err != nil {
		return err
	}
	return w.updateConfig(opts, nginxConfig.Json())
}

func (w *webServerConfigStore) Remove(ctx context.Context, opts metav1.WebServerOptions, pos metav1.ConfigContextPos) error {
	nginxConfig, err := w.getConfig(opts)
	if err != nil {
		return err
	}
	targetPos, err := w.parseContextPos(nginxConfig, pos)
	if err != nil {
		return errors.Errorf("failed to parse target position: %v", err)
	}
	father, targetIdx := targetPos.Position()
	if err = father.Remove(targetIdx).Error(); err != nil {
		return err
	}
	return w.updateConfig(opts, nginxConfig.Json())
}

func (w *webServerConfigStore) ModifyWithClone(_ context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	nginxConfig, err := w.getConfig(opts)
	if err != nil {
		return err
	}
	targetPos, err := w.parseContextPos(nginxConfig, ctxmeta.Position)
	if err != nil {
		return errors.Errorf("failed to parse target position: %v", err)
	}
	father, targetIdx := targetPos.Position()

	toBeClonedPos, err := w.parseContextPos(nginxConfig, ctxmeta.TargetContext.ConfigContextPos)
	if err != nil {
		return errors.Errorf("failed to parse context posision to be cloned: %v", err)
	}

	if err = father.Modify(toBeClonedPos.Target().Clone(), targetIdx).Error(); err != nil {
		return err
	}
	return w.updateConfig(opts, nginxConfig.Json())
}

func (w *webServerConfigStore) ModifyWithNew(_ context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	nginxConfig, err := w.getConfig(opts)
	if err != nil {
		return err
	}
	targetPos, err := w.parseContextPos(nginxConfig, ctxmeta.Position)
	if err != nil {
		return errors.Errorf("failed to parse target position: %v", err)
	}
	father, targetIdx := targetPos.Position()
	if err = father.Modify(newConfigContext(ctxmeta.TargetContext), targetIdx).Error(); err != nil {
		return err
	}
	return w.updateConfig(opts, nginxConfig.Json())
}

func (w *webServerConfigStore) ModifyContextValue(_ context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	nginxConfig, err := w.getConfig(opts)
	if err != nil {
		return err
	}
	targetPos, err := w.parseContextPos(nginxConfig, ctxmeta.Position)
	if err != nil {
		return errors.Errorf("failed to parse target position: %v", err)
	}
	if ctxmeta.TargetContext.ContextType != targetPos.Target().Type() {
		return errors.Errorf("the target context type for modification is inconsistent. the target context type: %s, the modified context type: %s.", targetPos.Target().Type(), ctxmeta.TargetContext.ContextType)
	}
	err = targetPos.Target().SetValue(ctxmeta.TargetContext.ContextValue)
	if err != nil {
		return errors.Errorf("failed to set value for target context: %v", err)
	}
	return w.updateConfig(opts, nginxConfig.Json())
}

func (w *webServerConfigStore) Move(_ context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	nginxConfig, err := w.getConfig(opts)
	if err != nil {
		return err
	}
	targetPos, err := w.parseContextPos(nginxConfig, ctxmeta.Position)
	if err != nil {
		return errors.Errorf("failed to parse target position: %v", err)
	}
	targetFather, targetIdx := targetPos.Position()
	toBeMovedPos, err := w.parseContextPos(nginxConfig, ctxmeta.TargetContext.ConfigContextPos)
	if err != nil {
		return errors.Errorf("failed to parse context posision to be moved: %v", err)
	}
	toBeMovedCtx := toBeMovedPos.Target()
	if err = toBeMovedCtx.Error(); err != nil {
		return errors.Errorf("failed to parse context posision to be moved: %v", err)
	}
	toBeMovedFather, toBeMovedIdx := toBeMovedPos.Position()
	maybeAfterTarget := toBeMovedIdx < targetIdx
	if !maybeAfterTarget {
		if err = toBeMovedFather.Remove(toBeMovedIdx).Error(); err != nil {
			return errors.Errorf("failed to removed context to be moved: %v", err)
		}
	}
	if err = targetFather.Insert(toBeMovedCtx, targetIdx).Error(); err != nil {
		return errors.Errorf("failed to insert context to be moved: %v", err)
	}
	if maybeAfterTarget {
		if err = toBeMovedFather.Remove(toBeMovedIdx).Error(); err != nil {
			return errors.Errorf("failed to removed context to be moved: %v", err)
		}
	}
	return w.updateConfig(opts, nginxConfig.Json())
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
		panic(errors.New("singleton web server config store is nil"))
	}
	return singletonWSCStore
}

func newConfigContext(meta metav1.NewConfigContextMeta) nginx_context.Context {
	ctx := local.NewContext(meta.ContextType, meta.ContextValue)
	for idx, childMeta := range meta.ChildrenContextMeta {
		ctx = ctx.Insert(newConfigContext(childMeta), idx)
	}
	return ctx
}
