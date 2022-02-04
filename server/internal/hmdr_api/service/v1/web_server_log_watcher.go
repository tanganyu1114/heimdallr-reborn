package v1

import (
	"context"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type WebServerLogWatcherSrv interface {
	Watch(ctx context.Context, opts metav1.WebServerLogOptions) (<-chan []byte, error)
}
