package bifrosts

import (
	"context"
	"gin-vue-admin/internal/pkg/bifrosts"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	bifrostapiv1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/marmotedu/errors"
	"sync"
)

var (
	wslwOnce           = new(sync.Once)
	singletonWSLWStore *webServerLogWatcherStore
)

type webServerLogWatcherStore struct {
	bm bifrosts.Manager
}

func (w *webServerLogWatcherStore) Watch(ctx context.Context, opts metav1.WebServerLogOptions) (<-chan []byte, context.CancelFunc, error) {
	bc, err := w.bm.GetBifrostClient(opts.WebServerOptions)
	if err != nil {
		return nil, nil, err
	}
	dataC, cancel, err := bc.WebServerLogWatcher().Watch(&bifrostapiv1.WebServerLogWatchRequest{
		ServerName:          &bifrostapiv1.ServerName{Name: opts.ServerName},
		LogName:             opts.LogName,
		FilteringRegexpRule: opts.FilteringRegexpRule,
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get watcher of server(%s)'s log(%s)", opts.ServerName, opts.LogName)
	}
	return dataC, cancel, nil
}

func newWebServerLogWatcherStore(store *bifrostsStore) *webServerLogWatcherStore {
	wslwOnce.Do(func() {
		if singletonWSLWStore == nil {
			singletonWSLWStore = &webServerLogWatcherStore{
				bm: store.bm,
			}
		}
	})
	if singletonWSLWStore == nil {
		panic(errors.New("singleton web server log watcher store is nil"))
	}
	return singletonWSLWStore
}
