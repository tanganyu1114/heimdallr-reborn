package v1

import (
	"context"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type WebServerBinCMDSrv interface {
	Exec(ctx context.Context, opts metav1.WebServerOptions, arg ...string) (metav1.WebServerBinCMDExecResponse, error)
}
