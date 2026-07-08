package logging

import (
	"context"
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/internal/hmdr_api/service/v1"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
	"go.uber.org/zap/zapcore"
)

type webServerBinCMDService struct {
	svc svcv1.Factory
}

var _ svcv1.WebServerBinCMDSrv = (*webServerBinCMDService)(nil)

func (w webServerBinCMDService) Exec(ctx context.Context, opts metav1.WebServerOptions, arg ...string) (resp metav1.WebServerBinCMDExecResponse, err error) {
	defer func() {
		level := zapcore.InfoLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "提交Web服务端二进制工具命令执行请求", map[string]any{"web-server-options": opts, "web-server-binary-cmd-arguments": arg}, resp, err)
	}()
	return w.svc.WebServerBinCMD().Exec(ctx, opts, arg...)
}

func newWebServerBinCMDService(svc svcv1.Factory) *webServerBinCMDService {
	return &webServerBinCMDService{svc: svc}
}
