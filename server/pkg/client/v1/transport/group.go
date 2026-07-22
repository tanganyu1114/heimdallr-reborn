package transport

import (
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	http_transport "github.com/go-kit/kit/transport/http"
)

// GroupTransport defines the interface for group related transport
type GroupTransport interface {
	// Get returns the get group client
	Get() httpclientv1.ClientBuilder[v1.IDOptions, modelclientv1.ResponseBody[*v1.Group]]
	// List returns the list groups client
	List() httpclientv1.ClientBuilder[v1.ListOptions, modelclientv1.ResponseBody[*v1.GroupList]]
}

// groupTransport implements GroupTransport interface
type groupTransport struct {
	getGroupClient   httpclientv1.ClientBuilder[v1.IDOptions, modelclientv1.ResponseBody[*v1.Group]]
	listGroupsClient httpclientv1.ClientBuilder[v1.ListOptions, modelclientv1.ResponseBody[*v1.GroupList]]
}

// newGroupTransport creates a new Groups transport
func newGroupTransport(transport *transport) GroupTransport {
	t := &groupTransport{
		getGroupClient: httpclientv1.NewClientBuilder[v1.IDOptions, modelclientv1.ResponseBody[*v1.Group]](
			httpclientv1.HTTPMethodGet,
			transport.baseURL+"/hmdrGroup/findHmdrGroup",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		listGroupsClient: httpclientv1.NewClientBuilder[v1.ListOptions, modelclientv1.ResponseBody[*v1.GroupList]](
			httpclientv1.HTTPMethodGet,
			transport.baseURL+"/hmdrGroup/getHmdrGroupList",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
	}
	return t
}

// Get returns the get group client
func (g *groupTransport) Get() httpclientv1.ClientBuilder[v1.IDOptions, modelclientv1.ResponseBody[*v1.Group]] {
	return g.getGroupClient
}

// List returns the list groups client
func (g *groupTransport) List() httpclientv1.ClientBuilder[v1.ListOptions, modelclientv1.ResponseBody[*v1.GroupList]] {
	return g.listGroupsClient
}
