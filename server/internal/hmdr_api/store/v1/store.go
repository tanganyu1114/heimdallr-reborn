package v1

type Factory interface {
	AgentInfos() AgentInfoStore
	Groups() GroupStore
	Hosts() HostStore
	WebServerConfigs() WebServerConfigStore
	WebServerLogWatchers() WebServerLogWatcherStore
	WebServerStatistics() WebServerStatisticsStore
	Close() error
}
