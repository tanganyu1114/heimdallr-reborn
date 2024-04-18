package fake

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/internal/pkg/bifrosts/fake"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
)

type WebServerConfigStore struct {
}

func newConfigContext(meta metav1.NewConfigContextMeta) (ctx nginx_context.Context) {
	switch meta.ContextType {
	case context_type.TypeDirective:
		ctx = local.NewDirective("", "")
		err := ctx.SetValue(meta.ContextValue)
		if err != nil {
			return nginx_context.ErrContext(err)
		}
	case context_type.TypeComment:
		ctx = local.NewComment(meta.ContextValue, false)
	case context_type.TypeInlineComment:
		ctx = local.NewComment(meta.ContextValue, true)
	default:
		ctx = local.NewContext(meta.ContextType, meta.ContextValue)
	}
	return
}

func (f WebServerConfigStore) GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error) {
	//TODO implement me
	panic("implement me")
}

func (f WebServerConfigStore) GetConfig(ctx context.Context, opts metav1.WebServerOptions) ([]string, error) {
	data, err := new(fake.ServiceClient).WebServerConfig().Get(opts.ServerName)
	if err != nil {
		return nil, err
	}
	c, err := configuration.NewNginxConfigFromJsonBytes(data)
	if err != nil {
		return nil, err
	}
	return c.Main().ConfigLines(false)
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
