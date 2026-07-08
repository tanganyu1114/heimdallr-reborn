package transport

import (
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	http_transport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/websocket"
)

// WebServerLogWatcherTransport defines the interface for web server log watcher related transport
type WebServerLogWatcherTransport interface {
	// Watch returns the watch logs client
	Watch() httpclientv1.Client[metav1.WebServerLogOptions, *websocket.Conn]
}

// webServerLogWatcherTransport implements WebServerLogWatcherTransport interface
type webServerLogWatcherTransport struct {
	watchClient httpclientv1.Client[metav1.WebServerLogOptions, *websocket.Conn]
}

// newWebServerLogWatcherTransport creates a new WebServerLogWatcher transport
func newWebServerLogWatcherTransport(transport *transport) WebServerLogWatcherTransport {
	t := &webServerLogWatcherTransport{
		watchClient: httpclientv1.NewHTTPClient[metav1.WebServerLogOptions, *websocket.Conn](
			httpclientv1.HTTPMethodGet,
			transport.baseURL+"/hmdrWebSocket/ws",
			http_transport.SetClient(transport.Client),
		),
	}
	return t
}

// Watch returns the watch logs client
func (w *webServerLogWatcherTransport) Watch() httpclientv1.Client[metav1.WebServerLogOptions, *websocket.Conn] {
	return w.watchClient
}
