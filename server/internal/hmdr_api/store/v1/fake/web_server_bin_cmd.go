package fake

import (
	"context"

	"github.com/marmotedu/errors"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/bifrosts/fake"
)

type WebServerBinCMDStore struct {
}

func (f WebServerBinCMDStore) Exec(ctx context.Context, opts metav1.WebServerOptions, arg ...string) (metav1.WebServerBinCMDExecResponse, error) {
	s, stdout, stderr, err := new(fake.ServiceClient).WebServerBinCMD().Exec(opts.ServerName, arg...)
	if err != nil {
		return metav1.WebServerBinCMDExecResponse{
			Successful: s,
			Stdout:     stdout,
			Stderr:     stderr,
		}, errors.Wrapf(err, "failed to submit request to server (%s) to execute command with arguments [%v]", opts.ServerName, arg)
	}
	return metav1.WebServerBinCMDExecResponse{
		Successful: s,
		Stdout:     stdout,
		Stderr:     stderr,
	}, nil
}
