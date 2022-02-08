package v1

import (
	"context"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type WebServerLogWatcherSrv interface {
	Watch(ctx context.Context, opts metav1.WebServerLogOptions) (<-chan []byte, context.CancelFunc, error)
}

type webServerLogWatcherService struct {
	store storev1.Factory
}

var _ WebServerLogWatcherSrv = (*webServerLogWatcherService)(nil)

func (w *webServerLogWatcherService) Watch(ctx context.Context, opts metav1.WebServerLogOptions) (<-chan []byte, context.CancelFunc, error) {
	return w.store.WebServerLogWatchers().Watch(ctx, opts)
}

func newWebServerLogWatchers(svc *service) WebServerLogWatcherSrv {
	return &webServerLogWatcherService{
		store: svc.store,
	}
}
