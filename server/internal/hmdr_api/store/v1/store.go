package v1

//go:generate mockgen -self_package=gin-vue-admin/internal/hmdr_api/store/v1 -destination=mock_store.go -package=v1 gin-vue-admin/internal/hmdr_api/store/v1 Factory,AgentInfoStore,GroupStore,HostStore,WebServerLogWatcherStore,WebServerStatisticsStore
//go:generate mockgen -source=web_server_config.go -destination=generics_mock_web_server_config.go -package=v1 gin-vue-admin/internal/hmdr_api/store/v1 WebServerConfigStore

type Factory interface {
	AgentInfos() AgentInfoStore
	Groups() GroupStore
	Hosts() HostStore
	WebServerConfigs() WebServerConfigStore
	WebServerLogWatchers() WebServerLogWatcherStore
	WebServerStatistics() WebServerStatisticsStore
	Close() error
}
