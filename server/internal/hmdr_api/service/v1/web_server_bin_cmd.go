package v1

import (
	"context"

	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
)

type WebServerBinCMDSrv interface {
	Exec(ctx context.Context, opts metav1.WebServerOptions, arg ...string) (metav1.WebServerBinCMDExecResponse, error)
}
