package v1

import (
	"context"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
)

type WebServerBinCMDStore interface {
	Exec(ctx context.Context, opts metav1.WebServerOptions, arg ...string) (metav1.WebServerBinCMDExecResponse, error)
}
