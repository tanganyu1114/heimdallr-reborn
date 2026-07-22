package transport

import (
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	http_transport "github.com/go-kit/kit/transport/http"
)

// HostTransport defines the interface for host related transport
type HostTransport interface {
	// Get returns the get host client
	Get() httpclientv1.ClientBuilder[v1.IDOptions, modelclientv1.ResponseBody[*v1.Host]]
	// List returns the list hosts client
	List() httpclientv1.ClientBuilder[v1.ListOptions, modelclientv1.ResponseBody[*v1.HostList]]
}

// hostTransport implements HostTransport interface
type hostTransport struct {
	getHostClient   httpclientv1.ClientBuilder[v1.IDOptions, modelclientv1.ResponseBody[*v1.Host]]
	listHostsClient httpclientv1.ClientBuilder[v1.ListOptions, modelclientv1.ResponseBody[*v1.HostList]]
}

// newHostTransport creates a new Hosts transport
func newHostTransport(transport *transport) HostTransport {
	t := &hostTransport{
		getHostClient: httpclientv1.NewClientBuilder[v1.IDOptions, modelclientv1.ResponseBody[*v1.Host]](
			httpclientv1.HTTPMethodGet,
			transport.baseURL+"/hmdrHost/findHmdrHost",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		listHostsClient: httpclientv1.NewClientBuilder[v1.ListOptions, modelclientv1.ResponseBody[*v1.HostList]](
			httpclientv1.HTTPMethodGet,
			transport.baseURL+"/hmdrHost/getHmdrHostList",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
	}
	return t
}

// Get returns the get host client
func (h *hostTransport) Get() httpclientv1.ClientBuilder[v1.IDOptions, modelclientv1.ResponseBody[*v1.Host]] {
	return h.getHostClient
}

// List returns the list hosts client
func (h *hostTransport) List() httpclientv1.ClientBuilder[v1.ListOptions, modelclientv1.ResponseBody[*v1.HostList]] {
	return h.listHostsClient
}
