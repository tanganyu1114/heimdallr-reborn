package fake

import (
	"context"
	storefake "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1/fake"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
)

type WebServerBinCMD struct {
}

func (f WebServerBinCMD) Exec(ctx context.Context, opts metav1.WebServerOptions, arg ...string) (metav1.WebServerBinCMDExecResponse, error) {
	return new(storefake.Store).WebServerBinCMD().Exec(ctx, opts, arg...)
}
