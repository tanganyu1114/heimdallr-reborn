package bifrosts

import (
	"context"
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/internal/hmdr_api/service/v1"
	storev1 "github.com/tanganyu1114/heimdallr-reborn/internal/hmdr_api/store/v1"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
)

type webServerLogWatcherService struct {
	store storev1.Factory
}

var _ svcv1.WebServerLogWatcherSrv = (*webServerLogWatcherService)(nil)

func (w *webServerLogWatcherService) Watch(ctx context.Context, opts metav1.WebServerLogOptions) (<-chan []byte, context.CancelFunc, error) {
	return w.store.WebServerLogWatchers().Watch(ctx, opts)
}

func newWebServerLogWatchers(svc *service) svcv1.WebServerLogWatcherSrv {
	return &webServerLogWatcherService{
		store: svc.store,
	}
}
