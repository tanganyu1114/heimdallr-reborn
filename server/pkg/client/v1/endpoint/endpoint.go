package endpoint

//go:generate mockgen -self_package=gin-vue-admin/pkg/client/v1/endpoint -destination=mock_endpoint.go -package=endpoint gin-vue-admin/pkg/client/v1/endpoint Factory,AgentInfoEndpoints,GroupEndpoints,HostEndpoints,WebServerConfigEndpoints,WebServerBinCMDEndpoints,WebServerStatisticsEndpoints
import (
	"sync"

	txpclientv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/transport"

	"github.com/marmotedu/errors"
)

// Factory defines the interface for all API endpoints
type Factory interface {
	// AgentInfos returns the Agent endpoint
	AgentInfos() AgentInfoEndpoints
	// Groups returns the Group endpoint
	Groups() GroupEndpoints
	// Hosts returns the Host endpoint
	Hosts() HostEndpoints
	// WebServerConfigs returns the WebServerConfig endpoint
	WebServerConfigs() WebServerConfigEndpoints
	// WebServerBinCMDs returns the WebServerBinCMD endpoint
	WebServerBinCMDs() WebServerBinCMDEndpoints
	// WebServerStatistics returns the WebServerStatistics endpoint
	WebServerStatistics() WebServerStatisticsEndpoints
	//// WebServerLogWatcher returns the WebServerLogWatcher endpoint
	//WebServerLogWatcher() WebServerLogWatcherEndpoints
}

// factory implements Factory interface
type factory struct {
	transport txpclientv1.Factory

	onceAgentInfo                sync.Once
	singletonAgentInfo           AgentInfoEndpoints
	onceGroup                    sync.Once
	singletonGroup               GroupEndpoints
	onceHost                     sync.Once
	singletonHost                HostEndpoints
	onceWebServerConfig          sync.Once
	singletonWebServerConfig     WebServerConfigEndpoints
	onceWebServerBinCMD          sync.Once
	singletonWebServerBinCMD     WebServerBinCMDEndpoints
	onceWebServerStatistics      sync.Once
	singletonWebServerStatistics WebServerStatisticsEndpoints
}

// AgentInfos returns the AgentInfos endpoint
func (f *factory) AgentInfos() AgentInfoEndpoints {
	f.onceAgentInfo.Do(func() {
		f.singletonAgentInfo = newAgentInfoEndpoints(f)
	})
	return f.singletonAgentInfo
}

// Groups returns the Group endpoint
func (f *factory) Groups() GroupEndpoints {
	f.onceGroup.Do(func() {
		f.singletonGroup = newGroupEndpoints(f)
	})
	return f.singletonGroup
}

// Hosts returns the Host endpoint
func (f *factory) Hosts() HostEndpoints {
	f.onceHost.Do(func() {
		f.singletonHost = newHostEndpoints(f)
	})
	return f.singletonHost
}

// WebServerConfigs returns the WebServerConfig endpoint
func (f *factory) WebServerConfigs() WebServerConfigEndpoints {
	f.onceWebServerConfig.Do(func() {
		f.singletonWebServerConfig = newWebServerConfigEndpoints(f)
	})
	return f.singletonWebServerConfig
}

// WebServerBinCMDs returns the WebServerBinCMD endpoint
func (f *factory) WebServerBinCMDs() WebServerBinCMDEndpoints {
	f.onceWebServerBinCMD.Do(func() {
		f.singletonWebServerBinCMD = newWebServerBinCMDEndpoints(f)
	})
	return f.singletonWebServerBinCMD
}

// WebServerStatistics returns the WebServerStatistics endpoint
func (f *factory) WebServerStatistics() WebServerStatisticsEndpoints {
	f.onceWebServerStatistics.Do(func() {
		f.singletonWebServerStatistics = newWebServerStatisticsEndpoints(f)
	})
	return f.singletonWebServerStatistics
}

//// WebServerLogWatcher returns the WebServerLogWatcher endpoint
//func (f *factory) WebServerLogWatcher() WebServerLogWatcherEndpoints {
//	return newWebServerLogWatcherEndpoints(f)
//}

// NewEndpoint creates a new endpoint factory
func NewEndpoint(transport txpclientv1.Factory) (Factory, error) {
	if transport == nil {
		return nil, errors.New("transport cannot be nil")
	}
	return &factory{
		transport: transport,
	}, nil
}
