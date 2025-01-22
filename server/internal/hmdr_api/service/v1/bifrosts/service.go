package bifrosts

import (
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
)

type service struct {
	store storev1.Factory
}

func NewService(store storev1.Factory) svcv1.Factory {
	return &service{
		store: store,
	}
}

func (s *service) AgentInfos() svcv1.AgentInfoSrv {
	return newAgentInfos(s)
}

func (s *service) Groups() svcv1.GroupSrv {
	return newGroups(s)
}

func (s *service) Hosts() svcv1.HostSrv {
	return newHosts(s)
}

func (s *service) WebServerConfigs() svcv1.WebServerConfigSrv {
	return newWebServerConfigs(s)
}

func (s *service) WebServerLogWatchers() svcv1.WebServerLogWatcherSrv {
	return newWebServerLogWatchers(s)
}

func (s *service) WebServerStatistics() svcv1.WebServerStatisticsSrv {
	return newWebServerStatistics(s)
}

func (s *service) WebServerBinCMD() svcv1.WebServerBinCMDSrv {
	return newWebServerBinCMD(s)
}
