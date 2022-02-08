package v1

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type WebServerConfigStore interface {
	GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error)
	GetConfig(ctx context.Context, opts metav1.WebServerOptions) ([]byte, error)
}
