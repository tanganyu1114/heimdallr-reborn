package fake

import (
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
)

type Store struct {
}

func (f Store) AgentInfos() storev1.AgentInfoStore {
	//TODO implement me
	panic("implement me")
}

func (f Store) Groups() storev1.GroupStore {
	//TODO implement me
	panic("implement me")
}

func (f Store) Hosts() storev1.HostStore {
	//TODO implement me
	panic("implement me")
}

func (f Store) WebServerConfigs() storev1.WebServerConfigStore {
	return new(WebServerConfigStore)
}

func (f Store) WebServerLogWatchers() storev1.WebServerLogWatcherStore {
	//TODO implement me
	panic("implement me")
}

func (f Store) WebServerStatistics() storev1.WebServerStatisticsStore {
	//TODO implement me
	panic("implement me")
}

func (f Store) Close() error {
	//TODO implement me
	panic("implement me")
}
