package transport

import (
	v1 "github.com/tanganyu1114/heimdallr-reborn/api/heimdallr_api/v1"

	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	http_transport "github.com/go-kit/kit/transport/http"
)

// AgentInfoTransport defines the interface for agent related transport
type AgentInfoTransport interface {
	// Get returns the get agent info client
	Get() httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.GroupInfo]]
}

// agentInfoTransport implements AgentInfoTransport interface
type agentInfoTransport struct {
	getAgentInfoClient httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.GroupInfo]]
}

// newAgentInfoTransport creates a new Agents transport
func newAgentInfoTransport(transport *transport) AgentInfoTransport {
	t := &agentInfoTransport{
		getAgentInfoClient: httpclientv1.NewClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.GroupInfo]](
			httpclientv1.HTTPMethodGet,
			transport.baseURL+"/agent/getAgentInfo",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
	}
	return t
}

// Get returns the get agent info client
func (a *agentInfoTransport) Get() httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.GroupInfo]] {
	return a.getAgentInfoClient
}
