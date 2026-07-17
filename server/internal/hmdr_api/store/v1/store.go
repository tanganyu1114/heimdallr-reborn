package v1

//go:generate mockgen -self_package=github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1 -destination=mock_store.go -package=v1 github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1 Factory,AgentInfoStore,GroupStore,HostStore,WebServerLogWatcherStore,WebServerStatisticsStore,WebServerConfigStore,WebServerBinCMDStore

type Factory interface {
	AgentInfos() AgentInfoStore
	Groups() GroupStore
	Hosts() HostStore
	WebServerConfigs() WebServerConfigStore
	WebServerLogWatchers() WebServerLogWatcherStore
	WebServerStatistics() WebServerStatisticsStore
	WebServerBinCMD() WebServerBinCMDStore
	Close() error
}
