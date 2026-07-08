package fake

import (
	"context"
	"github.com/marmotedu/errors"
	"github.com/tanganyu1114/heimdallr-reborn/internal/pkg/bifrosts/fake"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
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
