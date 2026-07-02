package endpoint

import (
	"gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	txpclientv1 "gin-vue-admin/pkg/client/v1/transport"
	"sync"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

// HostEndpoints defines the interface for host related endpoints
type HostEndpoints interface {
	// Get returns the get host endpoint
	Get() httpclientv1.Endpoint[metav1.IDOptions, *v1.Host]
	// List returns the list hosts endpoint
	List() httpclientv1.Endpoint[metav1.ListOptions, *v1.HostList]
}

// hostEndpoints implements HostEndpoints interface
type hostEndpoints struct {
	transport    txpclientv1.HostTransport
	onceGet      sync.Once
	getEndpoint  httpclientv1.Endpoint[metav1.IDOptions, *v1.Host]
	onceList     sync.Once
	listEndpoint httpclientv1.Endpoint[metav1.ListOptions, *v1.HostList]
}

// Get returns the get host endpoint
func (h *hostEndpoints) Get() httpclientv1.Endpoint[metav1.IDOptions, *v1.Host] {
	h.onceGet.Do(func() {
		h.getEndpoint = h.transport.Get().Build().Endpoint()
	})
	return h.getEndpoint
}

// List returns the list hosts endpoint
func (h *hostEndpoints) List() httpclientv1.Endpoint[metav1.ListOptions, *v1.HostList] {
	h.onceList.Do(func() {
		h.listEndpoint = h.transport.List().Build().Endpoint()
	})
	return h.listEndpoint
}

// newHostEndpoints creates a new Host endpoints
func newHostEndpoints(f *factory) HostEndpoints {
	return &hostEndpoints{
		transport: f.transport.Hosts(),
		onceGet:   sync.Once{},
		onceList:  sync.Once{},
	}
}
