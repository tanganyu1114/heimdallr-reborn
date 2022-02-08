package v1

type Factory interface {
	AgentInfos() AgentInfoStore
	Groups() GroupStore
	Hosts() HostStore
	WebServerConfigs() WebServerConfigStore
	WebServerLogWatchers() WebServerLogWatcherStore
	Close() error
}
