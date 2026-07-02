package endpoint

import (
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	txpclientv1 "gin-vue-admin/pkg/client/v1/transport"
	"sync"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

// WebServerBinCMDEndpoints defines the interface for web server binary command related endpoints
type WebServerBinCMDEndpoints interface {
	// Exec returns the exec endpoint
	Exec() httpclientv1.Endpoint[metav1.WebServerBinCMDExecRequest, *metav1.WebServerBinCMDExecResponse]
}

// webServerBinCMDEndpoints implements WebServerBinCMDEndpoints interface
type webServerBinCMDEndpoints struct {
	transport    txpclientv1.WebServerBinCMDTransport
	once         sync.Once
	execEndpoint httpclientv1.Endpoint[metav1.WebServerBinCMDExecRequest, *metav1.WebServerBinCMDExecResponse]
}

// Exec returns the exec endpoint
func (w *webServerBinCMDEndpoints) Exec() httpclientv1.Endpoint[metav1.WebServerBinCMDExecRequest, *metav1.WebServerBinCMDExecResponse] {
	w.once.Do(func() {
		w.execEndpoint = w.transport.Exec().Build().Endpoint()
	})
	return w.execEndpoint
}

// newWebServerBinCMDEndpoints creates a new WebServerBinCMD endpoints
func newWebServerBinCMDEndpoints(f *factory) WebServerBinCMDEndpoints {
	return &webServerBinCMDEndpoints{
		transport: f.transport.WebServerBinCMDs(),
		once:      sync.Once{},
	}
}
