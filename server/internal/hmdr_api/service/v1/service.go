package v1

import storev1 "gin-vue-admin/internal/hmdr_api/store/v1"

type Service interface {
	AgentInfos() AgentInfoSrv
	Groups() GroupSrv
	Hosts() HostSrv
	WebServerConfigs() WebServerConfigSrv
	WebServerLogWatchers() WebServerLogWatcherSrv
	WebServerStatistics() WebServerStatisticsSrv
}

type service struct {
	store storev1.Factory
}

func NewService(store storev1.Factory) Service {
	return &service{
		store: store,
	}
}

func (s *service) AgentInfos() AgentInfoSrv {
	return newAgentInfos(s)
}

func (s *service) Groups() GroupSrv {
	return newGroups(s)
}

func (s *service) Hosts() HostSrv {
	return newHosts(s)
}

func (s *service) WebServerConfigs() WebServerConfigSrv {
	return newWebServerConfigs(s)
}

func (s *service) WebServerLogWatchers() WebServerLogWatcherSrv {
	return newWebServerLogWatchers(s)
}

func (s *service) WebServerStatistics() WebServerStatisticsSrv {
	return newWebServerStatistics(s)
}
