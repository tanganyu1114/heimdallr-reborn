package fake

import (
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/service/v1"
)

type Service struct {
}

func (s Service) AgentInfos() svcv1.AgentInfoSrv {
	//TODO implement me
	panic("implement me")
}

func (s Service) Groups() svcv1.GroupSrv {
	//TODO implement me
	panic("implement me")
}

func (s Service) Hosts() svcv1.HostSrv {
	//TODO implement me
	panic("implement me")
}

func (s Service) WebServerConfigs() svcv1.WebServerConfigSrv {
	return new(WebServerConfigService)
}

func (s Service) WebServerLogWatchers() svcv1.WebServerLogWatcherSrv {
	//TODO implement me
	panic("implement me")
}

func (s Service) WebServerStatistics() svcv1.WebServerStatisticsSrv {
	//TODO implement me
	panic("implement me")
}

func (s Service) WebServerBinCMD() svcv1.WebServerBinCMDSrv {
	return new(WebServerBinCMD)
}
