package bifrosts

import (
	"fmt"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/global"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
	"gin-vue-admin/internal/pkg/bifrosts"
	"github.com/marmotedu/errors"
	"go.uber.org/zap"
	"sync"
)

var (
	storeOnce      = new(sync.Once)
	singletonStore *bifrostsStore
)

type bifrostsStore struct {
	bm bifrosts.Manager
}

func (b *bifrostsStore) AgentInfos() storev1.AgentInfoStore {
	return newAgentInfoStore(b)
}

func (b *bifrostsStore) Groups() storev1.GroupStore {
	return newGroupStore(b)
}

func (b *bifrostsStore) Hosts() storev1.HostStore {
	return newHostStore(b)
}

func (b *bifrostsStore) WebServerConfigs() storev1.WebServerConfigStore {
	return newWebServerConfigStore(b)
}

func (b *bifrostsStore) WebServerLogWatchers() storev1.WebServerLogWatcherStore {
	return newWebServerLogWatcherStore(b)
}

func (b *bifrostsStore) WebServerStatistics() storev1.WebServerStatisticsStore {
	return newWebServerStatisticsStore(b)
}

func (b *bifrostsStore) Close() error {
	return b.bm.RemoveAll()
}

func GetBifrostsStore() storev1.Factory {
	storeOnce.Do(func() {
		if singletonStore == nil {
			singletonStore = &bifrostsStore{
				bm: bifrosts.New(),
			}
			// init from global DB
			err := InitStoreFromGlobalDB()
			if err != nil {
				panic(err)
			}

		}
	})
	if singletonStore == nil {
		panic(errors.New("singleton bifrosts store is nil"))
	}
	return singletonStore
}

func InitStoreFromGlobalDB() error {
	if singletonStore == nil {
		return errors.New("singleton bifrosts store is nil")
	}
	err := singletonStore.Close()
	if err != nil {
		return errors.Wrap(err, "singleton bifrosts store init failed")
	}

	// init groups
	groupTable := global.GVA_DB.Model(&v1.Group{})
	var groups []v1.Group
	err = groupTable.Order("sequence").Find(&groups).Error
	if err != nil {
		global.GVA_LOG.Error("select hmdr_group error", zap.String("err", err.Error()))
		return err
	}

	for _, group := range groups {
		err := singletonStore.bm.InsertGroup(group)
		if err != nil {
			global.GVA_LOG.Error(fmt.Sprintf("init group(id: %d, name: %s) failed", group.ID, group.Name), zap.String("err", err.Error()))
			continue
		}

		// init hosts
		hostTable := global.GVA_DB.Model(&v1.Host{})
		var hosts []v1.Host
		err = hostTable.Where("group_id = ? and status = 1", group.ID).Order("sequence").Find(&hosts).Error
		if err != nil {
			global.GVA_LOG.Error(fmt.Sprintf("select hmdr_host error, where group_id = %d", group.ID), zap.String("err", err.Error()))
			continue
		}

		for _, host := range hosts {
			err = singletonStore.bm.InsertHost(host)
			if err != nil {
				global.GVA_LOG.Error(fmt.Sprintf("init host(id: %d, groupid: %d, groupname: %s) failed", host.ID, group.ID, group.Name), zap.String("err", err.Error()))
				continue
			}
		}

	}

	return nil

}
