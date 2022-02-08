package bifrosts

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/global"
	"gin-vue-admin/internal/pkg/bifrosts"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/pkg/sort_map"
	"github.com/marmotedu/errors"
	"go.uber.org/zap"
	"sync"
)

var (
	wscOnce           = new(sync.Once)
	singletonWSCStore *webServerConfigStore
)

type webServerConfigStore struct {
	bm bifrosts.Manager
}

func (w *webServerConfigStore) GetOptions(ctx context.Context) ([]v1.BifrostGroupMeta, error) {
	var result []v1.BifrostGroupMeta
	getOptsFromGroupFunc := func(keyer sort_map.Keyer, v interface{}) bool {
		group := v.(*bifrosts.Group)
		tmpGroup := v1.BifrostGroupMeta{
			ObjectMeta: metav1.ObjectMeta{
				Label: group.Meta.Name,
				Value: group.Meta.ID,
			},
			Children: make([]*v1.BifrostMeta, 0),
		}

		getOptFromBifrostFunc := func(keyer sort_map.Keyer, v interface{}) bool {
			bifrost := v.(*bifrosts.Bifrost)
			tmpHost := v1.BifrostMeta{
				ObjectMeta: metav1.ObjectMeta{
					Label: bifrost.Meta.Name,
					Value: bifrost.Meta.ID,
				},
				Children: make([]*v1.WebSrvMeta, 0),
			}

			metrics, err := bifrost.Client.WebServerStatus().Get()
			if err != nil {
				global.GVA_LOG.Error("failed to get web server status", zap.String("err", err.Error()))
				//return false  // 注意：如果报错，散列表后续元素将不被加载
			}
			for _, srvInfo := range metrics.StatusList {
				tmpHost.Children = append(tmpHost.Children,
					&v1.WebSrvMeta{
						ObjectMeta: metav1.ObjectMeta{
							Label: srvInfo.Name,
							Value: uint(srvInfo.Status),
						},
					},
				)
			}
			tmpGroup.Children = append(tmpGroup.Children, &tmpHost)
			return true
		}

		group.Bifrosts.Range(getOptFromBifrostFunc)
		result = append(result, tmpGroup)
		return true
	}

	w.bm.Range(getOptsFromGroupFunc)
	return result, nil
}

func (w webServerConfigStore) GetConfig(ctx context.Context, opts metav1.WebServerOptions) ([]byte, error) {
	bc, err := w.bm.GetBifrostClient(opts)
	if err != nil {
		return nil, err
	}
	data, err := bc.WebServerConfig().Get(opts.ServerName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get web server config")
	}
	return data, nil
}

func newWebServerConfigStore(store *bifrostsStore) *webServerConfigStore {
	wscOnce.Do(func() {
		if singletonWSCStore == nil {
			singletonWSCStore = &webServerConfigStore{
				bm: store.bm,
			}
		}
	})
	if singletonWSCStore == nil {
		panic(errors.New("singleton web server config store is nil"))
	}
	return singletonWSCStore
}
