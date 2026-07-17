package transport

//go:generate mockgen -self_package=github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/transport -destination=mock_transport.go -package=transport github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/transport Factory,SysUserTransport,AgentInfoTransport,GroupTransport,HostTransport,WebServerConfigTransport,WebServerBinCMDTransport,WebServerStatisticsTransport
import (
	"net/http"
	"sync"
	"time"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/marmotedu/errors"
)

// Factory defines the interface for creating transport clients
type Factory interface {
	// SysUsers returns the sysUser transport
	SysUsers() SysUserTransport
	// AgentInfos returns the Agent transport
	AgentInfos() AgentInfoTransport
	// Groups returns the Group transport
	Groups() GroupTransport
	// Hosts returns the Host transport
	Hosts() HostTransport
	// WebServerConfigs returns the WebServerConfig transport
	WebServerConfigs() WebServerConfigTransport
	// WebServerBinCMDs returns the WebServerBinCMD transport
	WebServerBinCMDs() WebServerBinCMDTransport
	// WebServerStatistics returns the WebServerStatistics transport
	WebServerStatistics() WebServerStatisticsTransport
	// TODO: WebServerLogWatchers transport is not implemented yet.
	// Implement when there is a requirement for websocket log watching.
	// WebServerLogWatchers returns the WebServerLogWatchers transport
	//WebServerLogWatchers() WebServerLogWatcherTransport
}

// transport implements Factory interface
type transport struct {
	Client  *http.Client
	baseURL string

	// SysUser transport
	onceSysUser      sync.Once
	singletonSysUser SysUserTransport

	// AgentInfo transport
	onceAgentInfo      sync.Once
	singletonAgentInfo AgentInfoTransport

	// Group transport
	onceGroup      sync.Once
	singletonGroup GroupTransport

	// Host transport
	onceHost      sync.Once
	singletonHost HostTransport

	// WebServerConfig transport
	onceWebServerConfig      sync.Once
	singletonWebServerConfig WebServerConfigTransport

	// WebServerBinCMD transport
	onceWebServerBinCMD      sync.Once
	singletonWebServerBinCMD WebServerBinCMDTransport

	// WebServerStatistics transport
	onceWebServerStatistics      sync.Once
	singletonWebServerStatistics WebServerStatisticsTransport
}

// NewTransport creates a new transport client
func NewTransport(client *http.Client, baseURL string) (Factory, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL cannot be empty")
	}
	if client == nil {
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return &transport{
		Client:  client,
		baseURL: baseURL,
	}, nil
}

// SysUsers returns the sysUser transport
func (t *transport) SysUsers() SysUserTransport {
	t.onceSysUser.Do(func() {
		if t.singletonSysUser == nil {
			t.singletonSysUser = newSysUserTransport(t)
		}
	})
	if t.singletonSysUser == nil {
		logV1.Fatal("sysUser transport client is nil")
		return nil
	}
	return t.singletonSysUser
}

// AgentInfos returns the Agent transport
func (t *transport) AgentInfos() AgentInfoTransport {
	t.onceAgentInfo.Do(func() {
		if t.singletonAgentInfo == nil {
			t.singletonAgentInfo = newAgentInfoTransport(t)
		}
	})
	if t.singletonAgentInfo == nil {
		logV1.Fatal("agentInfo transport client is nil")
		return nil
	}
	return t.singletonAgentInfo
}

// Groups returns the Group transport
func (t *transport) Groups() GroupTransport {
	t.onceGroup.Do(func() {
		if t.singletonGroup == nil {
			t.singletonGroup = newGroupTransport(t)
		}
	})
	if t.singletonGroup == nil {
		logV1.Fatal("group transport client is nil")
		return nil
	}
	return t.singletonGroup
}

// Hosts returns the Host transport
func (t *transport) Hosts() HostTransport {
	t.onceHost.Do(func() {
		if t.singletonHost == nil {
			t.singletonHost = newHostTransport(t)
		}
	})
	if t.singletonHost == nil {
		logV1.Fatal("host transport client is nil")
		return nil
	}
	return t.singletonHost
}

// WebServerConfigs returns the WebServerConfig transport
func (t *transport) WebServerConfigs() WebServerConfigTransport {
	t.onceWebServerConfig.Do(func() {
		if t.singletonWebServerConfig == nil {
			t.singletonWebServerConfig = newWebServerConfigTransport(t)
		}
	})
	if t.singletonWebServerConfig == nil {
		logV1.Fatal("webServerConfig transport client is nil")
		return nil
	}
	return t.singletonWebServerConfig
}

// WebServerBinCMDs returns the WebServerBinCMD transport
func (t *transport) WebServerBinCMDs() WebServerBinCMDTransport {
	t.onceWebServerBinCMD.Do(func() {
		if t.singletonWebServerBinCMD == nil {
			t.singletonWebServerBinCMD = newWebServerBinCMDTransport(t)
		}
	})
	if t.singletonWebServerBinCMD == nil {
		logV1.Fatal("webServerBinCMD transport client is nil")
		return nil
	}
	return t.singletonWebServerBinCMD
}

// WebServerStatistics returns the WebServerStatistics transport
func (t *transport) WebServerStatistics() WebServerStatisticsTransport {
	t.onceWebServerStatistics.Do(func() {
		if t.singletonWebServerStatistics == nil {
			t.singletonWebServerStatistics = newWebServerStatisticsTransport(t)
		}
	})
	if t.singletonWebServerStatistics == nil {
		logV1.Fatal("webServerStatistics transport client is nil")
		return nil
	}
	return t.singletonWebServerStatistics
}
