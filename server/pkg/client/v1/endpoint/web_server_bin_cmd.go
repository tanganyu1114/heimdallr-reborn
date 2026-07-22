package endpoint

import (
	"sync"

	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	txpclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/transport"

	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

// WebServerBinCMDEndpoints defines the interface for web server binary command related endpoints
type WebServerBinCMDEndpoints interface {
	// Exec returns the exec endpoint
	Exec() httpclientv1.Endpoint[metav1.WebServerBinCMDExecRequest, modelclientv1.ResponseBody[*metav1.WebServerBinCMDExecResponse]]
}

// webServerBinCMDEndpoints implements WebServerBinCMDEndpoints interface
type webServerBinCMDEndpoints struct {
	transport    txpclientv1.WebServerBinCMDTransport
	once         sync.Once
	execEndpoint httpclientv1.Endpoint[metav1.WebServerBinCMDExecRequest, modelclientv1.ResponseBody[*metav1.WebServerBinCMDExecResponse]]
}

// Exec returns the exec endpoint
func (w *webServerBinCMDEndpoints) Exec() httpclientv1.Endpoint[metav1.WebServerBinCMDExecRequest, modelclientv1.ResponseBody[*metav1.WebServerBinCMDExecResponse]] {
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
