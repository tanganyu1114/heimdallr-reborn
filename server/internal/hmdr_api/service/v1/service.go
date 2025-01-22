package v1

//go:generate mockgen -self_package=gin-vue-admin/internal/hmdr_api/service/v1 -destination=mock_service.go -package=v1 gin-vue-admin/internal/hmdr_api/service/v1 Factory,AgentInfoSrv,GroupSrv,HostSrv,WebServerLogWatcherSrv,WebServerStatisticsSrv,WebServerBinCMDSrv
//go:generate mockgen -source=web_server_config.go -destination=generics_mock_web_server_config.go -package=v1 gin-vue-admin/internal/hmdr_api/service/v1 WebServerConfigSrv

type Factory interface {
	AgentInfos() AgentInfoSrv
	Groups() GroupSrv
	Hosts() HostSrv
	WebServerConfigs() WebServerConfigSrv
	WebServerLogWatchers() WebServerLogWatcherSrv
	WebServerStatistics() WebServerStatisticsSrv
	WebServerBinCMD() WebServerBinCMDSrv
}
