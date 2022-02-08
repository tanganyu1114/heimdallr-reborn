package bifrosts

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/internal/pkg/bifrosts"
	"github.com/marmotedu/errors"
	"sync"
)

var (
	aiOnce           = new(sync.Once)
	singletonAIStore *agentInfoStore
)

type agentInfoStore struct {
	bm bifrosts.Manager
}

func (a agentInfoStore) Get(ctx context.Context) ([]v1.GroupInfo, error) {
	ginfos := a.bm.GetServersStatus()
	if ginfos == nil {
		return nil, errors.New("empty status")
	}
	return ginfos, nil
}

func (a *agentInfoStore) SyncAgentInfos() {
	a.bm.SyncServersStatus()
}

func newAgentInfoStore(store *bifrostsStore) *agentInfoStore {
	aiOnce.Do(func() {
		if singletonAIStore == nil {
			singletonAIStore = &agentInfoStore{
				bm: store.bm,
			}
		}
	})
	if singletonAIStore == nil {
		panic(errors.New("singleton agent info store is nil"))
	}
	return singletonAIStore
}
