package endpoint

import (
	"gin-vue-admin/api/heimdallr_api/v1"
	txpclientv1 "gin-vue-admin/pkg/client/v1/transport"
	"sync"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

// AgentInfoEndpoints defines the interface for agent info related endpoints
type AgentInfoEndpoints interface {
	// Get returns the get agent info endpoint
	Get() httpclientv1.Endpoint[httpclientv1.NilBody, []v1.GroupInfo]
}

// agentInfoEndpoints implements AgentInfoEndpoints interface
type agentInfoEndpoints struct {
	transport   txpclientv1.AgentInfoTransport
	once        sync.Once
	getEndpoint httpclientv1.Endpoint[httpclientv1.NilBody, []v1.GroupInfo]
}

// Get returns the get agent info endpoint
func (a *agentInfoEndpoints) Get() httpclientv1.Endpoint[httpclientv1.NilBody, []v1.GroupInfo] {
	a.once.Do(func() {
		a.getEndpoint = a.transport.Get().Build().Endpoint()
	})
	return a.getEndpoint
}

// newAgentInfoEndpoints creates a new AgentInfo endpoints
func newAgentInfoEndpoints(f *factory) AgentInfoEndpoints {
	return &agentInfoEndpoints{
		transport: f.transport.AgentInfos(),
		once:      sync.Once{},
	}
}
