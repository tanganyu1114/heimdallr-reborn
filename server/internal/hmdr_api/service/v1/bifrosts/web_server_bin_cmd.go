package bifrosts

import (
	"context"
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
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
