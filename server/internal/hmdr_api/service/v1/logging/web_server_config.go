package logging

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	"go.uber.org/zap/zapcore"
)

type webServerConfigService struct {
	svc svcv1.Factory
}

var _ svcv1.WebServerConfigSrv = (*webServerConfigService)(nil)

func (w *webServerConfigService) GetOptions(ctx context.Context) (meta []v1.BifrostGroupMeta, err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "查询Bifrost组元数据", nil, meta, err)
	}()
	return w.svc.WebServerConfigs().GetOptions(ctx)
}

func (w *webServerConfigService) GetConfig(ctx context.Context, opts metav1.WebServerOptions) (configmeta metav1.WebServerConfig, err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "查询服务配置", opts, configmeta, err)
	}()
	return w.svc.WebServerConfigs().GetConfig(ctx, opts)
}

func (w *webServerConfigService) GetContext(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos metav1.ConfigContextPos) (ngCtx nginx_context.Context, err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "查询目标配置上下文", map[string]any{"web-server-options": opts, "original-fingerprints": ofp, "config-context-position": pos}, ngCtx, err)
	}()
	return w.svc.WebServerConfigs().GetContext(ctx, opts, ofp, pos)
}

func (w *webServerConfigService) GetIncludedConfigs(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos metav1.ConfigContextPos) (includedConfigs []string, err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.WarnLevel
		}
		log(level, "查询目标`Include`上下文包含的配置文件路径列表", map[string]any{"web-server-options": opts, "original-fingerprints": ofp, "config-context-position": pos}, includedConfigs, err)
	}()
	return w.svc.WebServerConfigs().GetIncludedConfigs(ctx, opts, ofp, pos)
}

func (w *webServerConfigService) SearchContextPositions(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, kwmeta metav1.SearchKeywordsMeta) (poses []metav1.ConfigContextPos, err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.WarnLevel
		}
		log(level, "按关键字搜索匹配的上下文的坐标位置", map[string]any{"web-server-options": opts, "original-fingerprints": ofp, "search-keywords-meta": kwmeta}, poses, err)
	}()
	return w.svc.WebServerConfigs().SearchContextPositions(ctx, opts, ofp, kwmeta)
}

func (w *webServerConfigService) InsertWithClone(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "插入克隆的配置", map[string]any{"web-server-options": opts, "original-fingerprints": ofp, "config-context-meta": ctxmeta}, nil, err)
	}()
	return w.svc.WebServerConfigs().InsertWithClone(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigService) InsertWithNew(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "插入新增的配置", map[string]any{"web-server-options": opts, "original-fingerprints": ofp, "config-context-meta": ctxmeta}, nil, err)
	}()
	return w.svc.WebServerConfigs().InsertWithNew(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigService) Remove(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, pos metav1.ConfigContextPos) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "删除目标配置", map[string]any{"web-server-options": opts, "original-fingerprints": ofp, "config-context-position": pos}, nil, err)
	}()
	return w.svc.WebServerConfigs().Remove(ctx, opts, ofp, pos)
}

func (w *webServerConfigService) ModifyWithClone(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "修改为克隆的配置", map[string]any{"web-server-options": opts, "original-fingerprints": ofp, "config-context-meta": ctxmeta}, nil, err)
	}()
	return w.svc.WebServerConfigs().ModifyWithClone(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigService) ModifyWithNew(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "修改为新增的配置", map[string]any{"web-server-options": opts, "original-fingerprints": ofp, "config-context-meta": ctxmeta}, nil, err)
	}()
	return w.svc.WebServerConfigs().ModifyWithNew(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigService) ChangeContextEnabledState(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "启用/禁用目标上下文", map[string]any{"web-server-options": opts, "original-fingerprints": ofp, "config-context-meta": ctxmeta}, nil, err)
	}()
	return w.svc.WebServerConfigs().ChangeContextEnabledState(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigService) ModifyContextValue(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "修改目标上下文的参数值", map[string]any{"web-server-options": opts, "original-fingerprints": ofp, "config-context-meta": ctxmeta}, nil, err)
	}()
	return w.svc.WebServerConfigs().ModifyContextValue(ctx, opts, ofp, ctxmeta)
}

func (w *webServerConfigService) Move(ctx context.Context, opts metav1.WebServerOptions, ofp utilsV3.ConfigFingerprints, ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "移动目标配置", map[string]any{"web-server-options": opts, "original-fingerprints": ofp, "config-context-meta": ctxmeta}, nil, err)
	}()
	return w.svc.WebServerConfigs().Move(ctx, opts, ofp, ctxmeta)
}

func newWebServerConfigs(svc svcv1.Factory) *webServerConfigService {
	return &webServerConfigService{svc: svc}
}
