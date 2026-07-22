package v1

import (
	"context"

	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
)

type WebServerStatisticsStore interface {
	GetProxyServiceInfo(ctx context.Context, opts v1.WebServerOptions) ([]v1.ProxyServiceInfo, error)
	ConnectivityCheckOfProxyService(ctx context.Context, opts v1.WebServerOptions, proxyPassPos v1.ConfigContextPos) (v1.ProxyServiceInfo, error)
	ExportProxyServiceInfoToExcel(ctx context.Context, opts v1.WebServerOptions) ([]byte, error)
}
