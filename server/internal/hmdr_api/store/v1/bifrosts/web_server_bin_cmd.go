package bifrosts

import (
	"context"
	"github.com/marmotedu/errors"
	"github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/bifrosts"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
	"sync"
)

var (
	wsbcOnce           = new(sync.Once)
	singletonWSBCStore *webServerBinCMDStore
)

type webServerBinCMDStore struct {
	bm bifrosts.Manager
}

func (w *webServerBinCMDStore) Exec(ctx context.Context, opts metav1.WebServerOptions, arg ...string) (metav1.WebServerBinCMDExecResponse, error) {
	bc, err := w.bm.GetBifrostClient(opts)
	if err != nil {
		return metav1.WebServerBinCMDExecResponse{}, err
	}
	isSuccessful, stdout, stderr, err := bc.WebServerBinCMD().Exec(opts.ServerName, arg...)
	if err != nil {
		return metav1.WebServerBinCMDExecResponse{
			Successful: isSuccessful,
			Stdout:     stdout,
			Stderr:     stderr,
		}, errors.Wrapf(err, "failed to submit request to server (%s) to execute command with arguments [%v]", opts.ServerName, arg)
	}
	return metav1.WebServerBinCMDExecResponse{
		Successful: isSuccessful,
		Stdout:     stdout,
		Stderr:     stderr,
	}, nil
}

func newWebServerBinCMDStore(store *bifrostsStore) *webServerBinCMDStore {
	wsbcOnce.Do(func() {
		if singletonWSBCStore == nil {
			singletonWSBCStore = &webServerBinCMDStore{
				bm: store.bm,
			}
		}
	})
	if singletonWSBCStore == nil {
		panic(errors.New("the singleton web server binary command store is nil"))
	}
	return singletonWSBCStore
}
