package transport

import (
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	http_transport "github.com/go-kit/kit/transport/http"
)

// WebServerBinCMDTransport defines the interface for web server binary command related transport
type WebServerBinCMDTransport interface {
	// Exec returns the exec client
	Exec() httpclientv1.ClientBuilder[metav1.WebServerBinCMDExecRequest, modelclientv1.ResponseBody[*metav1.WebServerBinCMDExecResponse]]
}

// webServerBinCMDTransport implements WebServerBinCMDTransport interface
type webServerBinCMDTransport struct {
	execClient httpclientv1.ClientBuilder[metav1.WebServerBinCMDExecRequest, modelclientv1.ResponseBody[*metav1.WebServerBinCMDExecResponse]]
}

// newWebServerBinCMDTransport creates a new WebServerBinCMDs transport
func newWebServerBinCMDTransport(transport *transport) WebServerBinCMDTransport {
	t := &webServerBinCMDTransport{
		execClient: httpclientv1.NewClientBuilder[metav1.WebServerBinCMDExecRequest, modelclientv1.ResponseBody[*metav1.WebServerBinCMDExecResponse]](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/bin-cmd/exec",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
	}
	return t
}

// Exec returns the exec client
func (w *webServerBinCMDTransport) Exec() httpclientv1.ClientBuilder[metav1.WebServerBinCMDExecRequest, modelclientv1.ResponseBody[*metav1.WebServerBinCMDExecResponse]] {
	return w.execClient
}
