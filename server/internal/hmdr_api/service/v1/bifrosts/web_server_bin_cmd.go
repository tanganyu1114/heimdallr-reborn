package bifrosts

import (
	"context"

	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/service/v1"
	storev1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1"
)

type webServerBinCMDService struct {
	store storev1.Factory
}

var _ svcv1.WebServerBinCMDSrv = (*webServerBinCMDService)(nil)

func (w *webServerBinCMDService) Exec(ctx context.Context, opts metav1.WebServerOptions, arg ...string) (metav1.WebServerBinCMDExecResponse, error) {
	return w.store.WebServerBinCMD().Exec(ctx, opts, arg...)
}

func newWebServerBinCMD(svc *service) svcv1.WebServerBinCMDSrv {
	return &webServerBinCMDService{
		store: svc.store,
	}
}
