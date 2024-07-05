package fake

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/internal/pkg/bifrosts/fake"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/marmotedu/errors"
)

type WebServerConfigStore struct {
}

func newConfigContext(meta metav1.NewConfigContextMeta) nginx_context.Context {
	return local.NewContext(meta.ContextType, meta.ContextValue)
}

func (f WebServerConfigStore) GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error) {
	//TODO implement me
	panic("implement me")
}

func (f WebServerConfigStore) GetConfig(ctx context.Context, opts metav1.WebServerOptions) (configuration.NginxConfig, error) {
	data, err := new(fake.ServiceClient).WebServerConfig().Get(opts.ServerName)
	if err != nil {
		return nil, err
	}
	return configuration.NewNginxConfigFromJsonBytes(data)
}

func (f WebServerConfigStore) InsertWithClone(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	if len(ctxmeta.TargetContext.ContextPosPath) == 0 {
		return errors.Errorf("failed to parse context posision to be cloned: %v", errors.New("nginx config context pos path is null"))
	}
	return nil
}

func (f WebServerConfigStore) InsertWithNew(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	if err := newConfigContext(ctxmeta.TargetContext).Error(); err != nil {
		return err
	}
	return nil
}

func (f WebServerConfigStore) Remove(ctx context.Context, opts metav1.WebServerOptions, pos metav1.ConfigContextPos) error {
	if len(pos.ContextPosPath) == 0 {
		return errors.Errorf("failed to parse target position: %v", errors.New("nginx config context pos path is null"))
	}
	return nil
}

func (f WebServerConfigStore) ModifyWithClone(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	if len(ctxmeta.TargetContext.ContextPosPath) == 0 {
		return errors.Errorf("failed to parse context posision to be cloned: %v", errors.New("nginx config context pos path is null"))
	}
	return nil
}

func (f WebServerConfigStore) ModifyWithNew(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) error {
	if err := newConfigContext(ctxmeta.TargetContext).Error(); err != nil {
		return err
	}
	return nil
}

func (f WebServerConfigStore) parseContextPos(nginxconfig configuration.NginxConfig, pos metav1.ConfigContextPos) (nginx_context.Pos, error) {
	if len(pos.ContextPosPath) == 0 {
		return nginx_context.NullPos(), errors.New("nginx config context pos path is null")
	}

	nextFather := nginx_context.NullContext()
	nextFather, err := nginxconfig.Main().GetConfig(pos.Config)
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

func (f WebServerConfigStore) Move(ctx context.Context, opts metav1.WebServerOptions, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) error {
	ngconf, err := f.GetConfig(ctx, opts)
	if err != nil {
		return err
	}
	targetPos, err := f.parseContextPos(ngconf, ctxmeta.Position)
	if err != nil {
		return errors.Errorf("failed to parse target position: %v", err)
	}
	targetFather, targetIdx := targetPos.Position()
	toBeMovedPos, err := f.parseContextPos(ngconf, ctxmeta.TargetContext.ConfigContextPos)
	if err != nil {
		return errors.Errorf("failed to parse context posision to be moved: %v", err)
	}
	toBeMovedCtx := toBeMovedPos.Target()
	if err = toBeMovedCtx.Error(); err != nil {
		return errors.Errorf("failed to parse context posision to be moved: %v", err)
	}
	toBeMovedFather, toBeMovedIdx := toBeMovedPos.Position()
	if err = toBeMovedFather.Remove(toBeMovedIdx).Error(); err != nil {
		return errors.Errorf("failed to removed context to be moved: %v", err)
	}
	if err = targetFather.Insert(toBeMovedCtx, targetIdx).Error(); err != nil {
		return errors.Errorf("failed to insert context to be moved: %v", err)
	}
	return new(fake.ServiceClient).WebServerConfig().Update(opts.ServerName, ngconf.Json())
}
