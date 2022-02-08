package v1

import (
	"context"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type WebServerLogWatcherStore interface {
	Watch(ctx context.Context, opts metav1.WebServerLogOptions) (<-chan []byte, context.CancelFunc, error)
}
