package endpoint

import (
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
	txpclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/transport"

	"github.com/gorilla/websocket"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

// WebServerLogWatcherEndpoints defines the interface for web server log watcher related endpoints
type WebServerLogWatcherEndpoints interface {
	// Watch returns the watch logs endpoint
	Watch() httpclientv1.Endpoint[metav1.WebServerLogOptions, *websocket.Conn]
}

// webServerLogWatcherEndpoints implements WebServerLogWatcherEndpoints interface
type webServerLogWatcherEndpoints struct {
	transport txpclientv1.WebServerLogWatcherTransport
}

// Watch returns the watch logs endpoint
func (w *webServerLogWatcherEndpoints) Watch() httpclientv1.Endpoint[metav1.WebServerLogOptions, *websocket.Conn] {
	return w.transport.Watch().Endpoint()
}

//// newWebServerLogWatcherEndpoints creates a new WebServerLogWatcher endpoints
//func newWebServerLogWatcherEndpoints(f *factory) WebServerLogWatcherEndpoints {
//	return &webServerLogWatcherEndpoints{
//		transport: f.transport.WebServerLogWatchers(),
//	}
//}
