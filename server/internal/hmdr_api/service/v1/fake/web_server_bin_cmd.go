package fake

import (
	"context"
	storefake "gin-vue-admin/internal/hmdr_api/store/v1/fake"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type WebServerBinCMD struct {
}

func (f WebServerBinCMD) Exec(ctx context.Context, opts metav1.WebServerOptions, arg ...string) (metav1.WebServerBinCMDExecResponse, error) {
	return new(storefake.Store).WebServerBinCMD().Exec(ctx, opts, arg...)
}
