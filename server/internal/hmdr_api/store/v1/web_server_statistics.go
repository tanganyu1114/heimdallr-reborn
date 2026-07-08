package v1

import (
	"context"
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
)

type WebServerStatisticsStore interface {
	GetProxyServiceInfo(ctx context.Context, opts metav1.WebServerOptions) ([]v1.ProxyServiceInfo, error)
	ConnectivityCheckOfProxyService(ctx context.Context, opts metav1.WebServerOptions, proxyPassPos metav1.ConfigContextPos) (v1.ProxyServiceInfo, error)
	ExportProxyServiceInfoToExcel(ctx context.Context, opts metav1.WebServerOptions) ([]byte, error)
}
