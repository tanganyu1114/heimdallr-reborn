package endpoint

import (
	"gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	txpclientv1 "gin-vue-admin/pkg/client/v1/transport"
	"sync"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

// GroupEndpoints defines the interface for group related endpoints
type GroupEndpoints interface {
	// Get returns the get group endpoint
	Get() httpclientv1.Endpoint[metav1.IDOptions, *v1.Group]
	// List returns the list groups endpoint
	List() httpclientv1.Endpoint[metav1.ListOptions, *v1.GroupList]
}

// groupEndpoints implements GroupEndpoints interface
type groupEndpoints struct {
	transport    txpclientv1.GroupTransport
	onceGet      sync.Once
	getEndpoint  httpclientv1.Endpoint[metav1.IDOptions, *v1.Group]
	onceList     sync.Once
	listEndpoint httpclientv1.Endpoint[metav1.ListOptions, *v1.GroupList]
}

// Get returns the get group endpoint
func (g *groupEndpoints) Get() httpclientv1.Endpoint[metav1.IDOptions, *v1.Group] {
	g.onceGet.Do(func() {
		g.getEndpoint = g.transport.Get().Build().Endpoint()
	})
	return g.getEndpoint
}

// List returns the list groups endpoint
func (g *groupEndpoints) List() httpclientv1.Endpoint[metav1.ListOptions, *v1.GroupList] {
	g.onceList.Do(func() {
		g.listEndpoint = g.transport.List().Build().Endpoint()
	})
	return g.listEndpoint
}

// newGroupEndpoints creates a new Group endpoints
func newGroupEndpoints(f *factory) GroupEndpoints {
	return &groupEndpoints{
		transport: f.transport.Groups(),
		onceGet:   sync.Once{},
		onceList:  sync.Once{},
	}
}
