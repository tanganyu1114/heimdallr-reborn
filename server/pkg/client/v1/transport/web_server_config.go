package transport

import (
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	http_transport "github.com/go-kit/kit/transport/http"
)

// WebServerConfigTransport defines the interface for web server config related transport
type WebServerConfigTransport interface {
	// GetOptions returns the get options client
	GetOptions() httpclientv1.ClientBuilder[httpclientv1.NilBody, []v1.BifrostGroupMeta]
	// GetConfigTextLines returns the get config text lines client
	GetConfigTextLines() httpclientv1.ClientBuilder[metav1.WebServerOptions, string]
	// GetContextTextLines returns the get context text lines client
	GetContextTextLines() httpclientv1.ClientBuilder[metav1.WebServerConfigTargetContextOptions, string]
	// GetConfig returns the get config client
	GetConfig() httpclientv1.ClientBuilder[metav1.WebServerOptions, *metav1.WebServerConfig]
	// GetIncludedConfigs returns the get included configs client
	GetIncludedConfigs() httpclientv1.ClientBuilder[metav1.WebServerConfigTargetContextOptions, []string]
	// SearchContextPositions returns the search context positions client
	SearchContextPositions() httpclientv1.ClientBuilder[metav1.WebServerConfigContextPosSearchOptions, []metav1.ConfigContextPos]
	// InsertWithClone returns the insert with clone client
	InsertWithClone() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody]
	// InsertWithNew returns the insert with new client
	InsertWithNew() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody]
	// Remove returns the remove client
	Remove() httpclientv1.ClientBuilder[metav1.WebServerConfigTargetContextOptions, httpclientv1.NilBody]
	// ModifyContextValue returns the modify context value client
	ModifyContextValue() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody]
	// ModifyWithClone returns the modify with clone client
	ModifyWithClone() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody]
	// ChangeContextEnabledState returns the change context enabled state client
	ChangeContextEnabledState() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], httpclientv1.NilBody]
	// ModifyWithNew returns the modify with new client
	ModifyWithNew() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody]
	// Move returns the move client
	Move() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody]
}

// webServerConfigTransport implements WebServerConfigTransport interface
type webServerConfigTransport struct {
	getOptionsClient                httpclientv1.ClientBuilder[httpclientv1.NilBody, []v1.BifrostGroupMeta]
	getConfigTextLinesClient        httpclientv1.ClientBuilder[metav1.WebServerOptions, string]
	getContextTextLinesClient       httpclientv1.ClientBuilder[metav1.WebServerConfigTargetContextOptions, string]
	getConfigClient                 httpclientv1.ClientBuilder[metav1.WebServerOptions, *metav1.WebServerConfig]
	getIncludedConfigsClient        httpclientv1.ClientBuilder[metav1.WebServerConfigTargetContextOptions, []string]
	searchContextPositionsClient    httpclientv1.ClientBuilder[metav1.WebServerConfigContextPosSearchOptions, []metav1.ConfigContextPos]
	insertWithCloneClient           httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody]
	insertWithNewClient             httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody]
	removeClient                    httpclientv1.ClientBuilder[metav1.WebServerConfigTargetContextOptions, httpclientv1.NilBody]
	modifyContextValueClient        httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody]
	modifyWithCloneClient           httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody]
	changeContextEnabledStateClient httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], httpclientv1.NilBody]
	modifyWithNewClient             httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody]
	moveClient                      httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody]
}

// newWebServerConfigTransport creates a new WebServerConfigs transport
func newWebServerConfigTransport(transport *transport) WebServerConfigTransport {
	t := &webServerConfigTransport{
		getOptionsClient: httpclientv1.NewClientBuilder[httpclientv1.NilBody, []v1.BifrostGroupMeta](
			httpclientv1.HTTPMethodGet,
			transport.baseURL+"/conf/getOptions",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		getConfigTextLinesClient: httpclientv1.NewClientBuilder[metav1.WebServerOptions, string](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/conf/getConfInfo",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		getContextTextLinesClient: httpclientv1.NewClientBuilder[metav1.WebServerConfigTargetContextOptions, string](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/conf/get-context-text",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		getConfigClient: httpclientv1.NewClientBuilder[metav1.WebServerOptions, *metav1.WebServerConfig](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/conf/get-conf-struct",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		getIncludedConfigsClient: httpclientv1.NewClientBuilder[metav1.WebServerConfigTargetContextOptions, []string](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/conf/get-includes",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		searchContextPositionsClient: httpclientv1.NewClientBuilder[metav1.WebServerConfigContextPosSearchOptions, []metav1.ConfigContextPos](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/conf/search-ctx-poses",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		insertWithCloneClient: httpclientv1.NewClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/conf/insert-clone-ctx",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		insertWithNewClient: httpclientv1.NewClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/conf/insert-new-ctx",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		removeClient: httpclientv1.NewClientBuilder[metav1.WebServerConfigTargetContextOptions, httpclientv1.NilBody](
			httpclientv1.HTTPMethodDelete,
			transport.baseURL+"/conf/remove-ctx",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		modifyContextValueClient: httpclientv1.NewClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/conf/modify-ctx-value",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		modifyWithCloneClient: httpclientv1.NewClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/conf/modify-clone-ctx",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		changeContextEnabledStateClient: httpclientv1.NewClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], httpclientv1.NilBody](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/conf/change-ctx-enabled-state",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		modifyWithNewClient: httpclientv1.NewClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/conf/modify-new-ctx",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		moveClient: httpclientv1.NewClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/conf/move-ctx",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
	}
	return t
}

// GetOptions returns the get options client
func (w *webServerConfigTransport) GetOptions() httpclientv1.ClientBuilder[httpclientv1.NilBody, []v1.BifrostGroupMeta] {
	return w.getOptionsClient
}

// GetConfigTextLines returns the get config text lines client
func (w *webServerConfigTransport) GetConfigTextLines() httpclientv1.ClientBuilder[metav1.WebServerOptions, string] {
	return w.getConfigTextLinesClient
}

// GetContextTextLines returns the get context text lines client
func (w *webServerConfigTransport) GetContextTextLines() httpclientv1.ClientBuilder[metav1.WebServerConfigTargetContextOptions, string] {
	return w.getContextTextLinesClient
}

// GetConfig returns the get config client
func (w *webServerConfigTransport) GetConfig() httpclientv1.ClientBuilder[metav1.WebServerOptions, *metav1.WebServerConfig] {
	return w.getConfigClient
}

// GetIncludedConfigs returns the get included configs client
func (w *webServerConfigTransport) GetIncludedConfigs() httpclientv1.ClientBuilder[metav1.WebServerConfigTargetContextOptions, []string] {
	return w.getIncludedConfigsClient
}

// SearchContextPositions returns the search context positions client
func (w *webServerConfigTransport) SearchContextPositions() httpclientv1.ClientBuilder[metav1.WebServerConfigContextPosSearchOptions, []metav1.ConfigContextPos] {
	return w.searchContextPositionsClient
}

// InsertWithClone returns the insert with clone client
func (w *webServerConfigTransport) InsertWithClone() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody] {
	return w.insertWithCloneClient
}

// InsertWithNew returns the insert with new client
func (w *webServerConfigTransport) InsertWithNew() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody] {
	return w.insertWithNewClient
}

// Remove returns the remove client
func (w *webServerConfigTransport) Remove() httpclientv1.ClientBuilder[metav1.WebServerConfigTargetContextOptions, httpclientv1.NilBody] {
	return w.removeClient
}

// ModifyContextValue returns the modify context value client
func (w *webServerConfigTransport) ModifyContextValue() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody] {
	return w.modifyContextValueClient
}

// ModifyWithClone returns the modify with clone client
func (w *webServerConfigTransport) ModifyWithClone() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody] {
	return w.modifyWithCloneClient
}

// ChangeContextEnabledState returns the change context enabled state client
func (w *webServerConfigTransport) ChangeContextEnabledState() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], httpclientv1.NilBody] {
	return w.changeContextEnabledStateClient
}

// ModifyWithNew returns the modify with new client
func (w *webServerConfigTransport) ModifyWithNew() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody] {
	return w.modifyWithNewClient
}

// Move returns the move client
func (w *webServerConfigTransport) Move() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody] {
	return w.moveClient
}
