package service

//go:generate mockgen -self_package=gin-vue-admin/pkg/client/v1/service -destination=mock_service.go -package=service gin-vue-admin/pkg/client/v1/service Factory,AgentInfoService,GroupService,HostService,WebServerConfigService,WebServerBinCMDService,WebServerStatisticsService
import (
	"context"
	"sync"

	epclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/endpoint"

	"github.com/marmotedu/errors"
)

// Factory defines the service layer interface
type Factory interface {
	// AgentInfos returns the AgentInfo service
	AgentInfos() AgentInfoService
	// Groups returns the Group service
	Groups() GroupService
	// Hosts returns the Host service
	Hosts() HostService
	// WebServerConfigs returns the WebServerConfig service
	WebServerConfigs() WebServerConfigService
	// WebServerBinCMDs returns the WebServerBinCMD service
	WebServerBinCMDs() WebServerBinCMDService
	// WebServerStatistics returns the WebServerStatistics service
	WebServerStatistics() WebServerStatisticsService
	//// WebServerLogWatcher returns the WebServerLogWatcher service
	//WebServerLogWatcher() WebServerLogWatcherService
}

// factory implements Factory interface
type factory struct {
	ctx context.Context
	eps epclientv1.Factory

	onceAgentInfo                sync.Once
	singletonAgentInfo           AgentInfoService
	onceGroup                    sync.Once
	singletonGroup               GroupService
	onceHost                     sync.Once
	singletonHost                HostService
	onceWebServerConfig          sync.Once
	singletonWebServerConfig     WebServerConfigService
	onceWebServerBinCMD          sync.Once
	singletonWebServerBinCMD     WebServerBinCMDService
	onceWebServerStatistics      sync.Once
	singletonWebServerStatistics WebServerStatisticsService
}

// AgentInfos returns the AgentInfo service
func (f *factory) AgentInfos() AgentInfoService {
	f.onceAgentInfo.Do(func() {
		f.singletonAgentInfo = newAgentInfoService(f)
	})
	return f.singletonAgentInfo
}

// Groups returns the Group service
func (f *factory) Groups() GroupService {
	f.onceGroup.Do(func() {
		f.singletonGroup = newGroupService(f)
	})
	return f.singletonGroup
}

// Hosts returns the Host service
func (f *factory) Hosts() HostService {
	f.onceHost.Do(func() {
		f.singletonHost = newHostService(f)
	})
	return f.singletonHost
}

// WebServerConfigs returns the WebServerConfig service
func (f *factory) WebServerConfigs() WebServerConfigService {
	f.onceWebServerConfig.Do(func() {
		f.singletonWebServerConfig = newWebServerConfigService(f)
	})
	return f.singletonWebServerConfig
}

// WebServerBinCMDs returns the WebServerBinCMD service
func (f *factory) WebServerBinCMDs() WebServerBinCMDService {
	f.onceWebServerBinCMD.Do(func() {
		f.singletonWebServerBinCMD = newWebServerBinCMDService(f)
	})
	return f.singletonWebServerBinCMD
}

// WebServerStatistics returns the WebServerStatistics service
func (f *factory) WebServerStatistics() WebServerStatisticsService {
	f.onceWebServerStatistics.Do(func() {
		f.singletonWebServerStatistics = newWebServerStatisticsService(f)
	})
	return f.singletonWebServerStatistics
}

//// WebServerLogWatcher returns the WebServerLogWatcher service
//func (f *factory) WebServerLogWatcher() WebServerLogWatcherService {
//	return newWebServerLogWatcherService(f)
//}

// NewService creates a new service factory
func NewService(ctx context.Context, eps epclientv1.Factory) (Factory, error) {
	if eps == nil {
		return nil, errors.New("endpoint factory cannot be nil")
	}
	return &factory{ctx: ctx, eps: eps}, nil
}
