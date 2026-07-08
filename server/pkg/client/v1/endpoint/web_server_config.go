package endpoint

import (
	v1 "github.com/tanganyu1114/heimdallr-reborn/api/heimdallr_api/v1"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
	txpclientv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/transport"
	"sync"

	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

// WebServerConfigEndpoints defines the interface for web server config related endpoints
type WebServerConfigEndpoints interface {
	// GetOptions returns the get options endpoint
	GetOptions() httpclientv1.Endpoint[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.BifrostGroupMeta]]
	// GetConfigTextLines returns the get config text lines endpoint
	GetConfigTextLines() httpclientv1.Endpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[string]]
	// GetContextTextLines returns the get context text lines endpoint
	GetContextTextLines() httpclientv1.Endpoint[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[string]]
	// GetConfig returns the get config endpoint
	GetConfig() httpclientv1.Endpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[*modelclientv1.WebServerConfig]]
	// GetIncludedConfigs returns the get included configs endpoint
	GetIncludedConfigs() httpclientv1.Endpoint[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[[]string]]
	// SearchContextPositions returns the search context positions endpoint
	SearchContextPositions() httpclientv1.Endpoint[metav1.WebServerConfigContextPosSearchOptions, modelclientv1.ResponseBody[[]metav1.ConfigContextPos]]
	// InsertWithClone returns the insert with clone endpoint
	InsertWithClone() httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	// InsertWithNew returns the insert with new endpoint
	InsertWithNew() httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	// Remove returns the remove endpoint
	Remove() httpclientv1.Endpoint[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	// UpdateConfig returns the update config endpoint
	UpdateConfig() httpclientv1.Endpoint[*metav1.WebServerConfigUpdateOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	// ModifyContextValue returns the modify context value endpoint
	ModifyContextValue() httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	// ModifyWithClone returns the modify with clone endpoint
	ModifyWithClone() httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	// ChangeContextEnabledState returns the change context enabled state endpoint
	ChangeContextEnabledState() httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	// ModifyWithNew returns the modify with new endpoint
	ModifyWithNew() httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	// Move returns the move endpoint
	Move() httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
}

// webServerConfigEndpoints implements WebServerConfigEndpoints interface
type webServerConfigEndpoints struct {
	transport                         txpclientv1.WebServerConfigTransport
	onceGetOptions                    sync.Once
	getOptionsEndpoint                httpclientv1.Endpoint[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.BifrostGroupMeta]]
	onceGetConfigTextLines            sync.Once
	getConfigTextLinesEndpoint        httpclientv1.Endpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[string]]
	onceGetContextTextLines           sync.Once
	getContextTextLinesEndpoint       httpclientv1.Endpoint[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[string]]
	onceGetConfig                     sync.Once
	getConfigEndpoint                 httpclientv1.Endpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[*modelclientv1.WebServerConfig]]
	onceGetIncludedConfigs            sync.Once
	getIncludedConfigsEndpoint        httpclientv1.Endpoint[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[[]string]]
	onceSearchContextPositions        sync.Once
	searchContextPositionsEndpoint    httpclientv1.Endpoint[metav1.WebServerConfigContextPosSearchOptions, modelclientv1.ResponseBody[[]metav1.ConfigContextPos]]
	onceInsertWithClone               sync.Once
	insertWithCloneEndpoint           httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	onceInsertWithNew                 sync.Once
	insertWithNewEndpoint             httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	onceRemove                        sync.Once
	removeEndpoint                    httpclientv1.Endpoint[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	onceUpdateConfig                  sync.Once
	updateConfigEndpoint              httpclientv1.Endpoint[*metav1.WebServerConfigUpdateOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	onceModifyContextValue            sync.Once
	modifyContextValueEndpoint        httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	onceModifyWithClone               sync.Once
	modifyWithCloneEndpoint           httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	onceChangeContextEnabledState     sync.Once
	changeContextEnabledStateEndpoint httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	onceModifyWithNew                 sync.Once
	modifyWithNewEndpoint             httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	onceMove                          sync.Once
	moveEndpoint                      httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
}

// GetOptions returns the get options endpoint
func (w *webServerConfigEndpoints) GetOptions() httpclientv1.Endpoint[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.BifrostGroupMeta]] {
	w.onceGetOptions.Do(func() {
		w.getOptionsEndpoint = w.transport.GetOptions().Build().Endpoint()
	})
	return w.getOptionsEndpoint
}

// GetConfigTextLines returns the get config text lines endpoint
func (w *webServerConfigEndpoints) GetConfigTextLines() httpclientv1.Endpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[string]] {
	w.onceGetConfigTextLines.Do(func() {
		w.getConfigTextLinesEndpoint = w.transport.GetConfigTextLines().Build().Endpoint()
	})
	return w.getConfigTextLinesEndpoint
}

// GetContextTextLines returns the get context text lines endpoint
func (w *webServerConfigEndpoints) GetContextTextLines() httpclientv1.Endpoint[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[string]] {
	w.onceGetContextTextLines.Do(func() {
		w.getContextTextLinesEndpoint = w.transport.GetContextTextLines().Build().Endpoint()
	})
	return w.getContextTextLinesEndpoint
}

// GetConfig returns the get config endpoint
func (w *webServerConfigEndpoints) GetConfig() httpclientv1.Endpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[*modelclientv1.WebServerConfig]] {
	w.onceGetConfig.Do(func() {
		w.getConfigEndpoint = w.transport.GetConfig().Build().Endpoint()
	})
	return w.getConfigEndpoint
}

// GetIncludedConfigs returns the get included configs endpoint
func (w *webServerConfigEndpoints) GetIncludedConfigs() httpclientv1.Endpoint[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[[]string]] {
	w.onceGetIncludedConfigs.Do(func() {
		w.getIncludedConfigsEndpoint = w.transport.GetIncludedConfigs().Build().Endpoint()
	})
	return w.getIncludedConfigsEndpoint
}

// SearchContextPositions returns the search context positions endpoint
func (w *webServerConfigEndpoints) SearchContextPositions() httpclientv1.Endpoint[metav1.WebServerConfigContextPosSearchOptions, modelclientv1.ResponseBody[[]metav1.ConfigContextPos]] {
	w.onceSearchContextPositions.Do(func() {
		w.searchContextPositionsEndpoint = w.transport.SearchContextPositions().Build().Endpoint()
	})
	return w.searchContextPositionsEndpoint
}

// InsertWithClone returns the insert with clone endpoint
func (w *webServerConfigEndpoints) InsertWithClone() httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	w.onceInsertWithClone.Do(func() {
		w.insertWithCloneEndpoint = w.transport.InsertWithClone().Build().Endpoint()
	})
	return w.insertWithCloneEndpoint
}

// InsertWithNew returns the insert with new endpoint
func (w *webServerConfigEndpoints) InsertWithNew() httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	w.onceInsertWithNew.Do(func() {
		w.insertWithNewEndpoint = w.transport.InsertWithNew().Build().Endpoint()
	})
	return w.insertWithNewEndpoint
}

// Remove returns the remove endpoint
func (w *webServerConfigEndpoints) Remove() httpclientv1.Endpoint[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	w.onceRemove.Do(func() {
		w.removeEndpoint = w.transport.Remove().Build().Endpoint()
	})
	return w.removeEndpoint
}

// UpdateConfig returns the update config endpoint
func (w *webServerConfigEndpoints) UpdateConfig() httpclientv1.Endpoint[*metav1.WebServerConfigUpdateOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	w.onceUpdateConfig.Do(func() {
		w.updateConfigEndpoint = w.transport.UpdateConfig().Build().Endpoint()
	})
	return w.updateConfigEndpoint
}

// ModifyContextValue returns the modify context value endpoint
func (w *webServerConfigEndpoints) ModifyContextValue() httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	w.onceModifyContextValue.Do(func() {
		w.modifyContextValueEndpoint = w.transport.ModifyContextValue().Build().Endpoint()
	})
	return w.modifyContextValueEndpoint
}

// ModifyWithClone returns the modify with clone endpoint
func (w *webServerConfigEndpoints) ModifyWithClone() httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	w.onceModifyWithClone.Do(func() {
		w.modifyWithCloneEndpoint = w.transport.ModifyWithClone().Build().Endpoint()
	})
	return w.modifyWithCloneEndpoint
}

// ChangeContextEnabledState returns the change context enabled state endpoint
func (w *webServerConfigEndpoints) ChangeContextEnabledState() httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	w.onceChangeContextEnabledState.Do(func() {
		w.changeContextEnabledStateEndpoint = w.transport.ChangeContextEnabledState().Build().Endpoint()
	})
	return w.changeContextEnabledStateEndpoint
}

// ModifyWithNew returns the modify with new endpoint
func (w *webServerConfigEndpoints) ModifyWithNew() httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	w.onceModifyWithNew.Do(func() {
		w.modifyWithNewEndpoint = w.transport.ModifyWithNew().Build().Endpoint()
	})
	return w.modifyWithNewEndpoint
}

// Move returns the move endpoint
func (w *webServerConfigEndpoints) Move() httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	w.onceMove.Do(func() {
		w.moveEndpoint = w.transport.Move().Build().Endpoint()
	})
	return w.moveEndpoint
}

// newWebServerConfigEndpoints creates a new WebServerConfig endpoints
func newWebServerConfigEndpoints(f *factory) WebServerConfigEndpoints {
	return &webServerConfigEndpoints{
		transport:                     f.transport.WebServerConfigs(),
		onceGetOptions:                sync.Once{},
		onceGetConfigTextLines:        sync.Once{},
		onceGetContextTextLines:       sync.Once{},
		onceGetConfig:                 sync.Once{},
		onceGetIncludedConfigs:        sync.Once{},
		onceSearchContextPositions:    sync.Once{},
		onceInsertWithClone:           sync.Once{},
		onceInsertWithNew:             sync.Once{},
		onceRemove:                    sync.Once{},
		onceModifyContextValue:        sync.Once{},
		onceModifyWithClone:           sync.Once{},
		onceChangeContextEnabledState: sync.Once{},
		onceModifyWithNew:             sync.Once{},
		onceMove:                      sync.Once{},
	}
}
