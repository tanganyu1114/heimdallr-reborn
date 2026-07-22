package fake

import (
	"context"

	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	storefake "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1/fake"
)

type WebServerBinCMD struct {
}

func (f WebServerBinCMD) Exec(ctx context.Context, opts metav1.WebServerOptions, arg ...string) (metav1.WebServerBinCMDExecResponse, error) {
	return new(storefake.Store).WebServerBinCMD().Exec(ctx, opts, arg...)
}
