package auth

import (
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

type webServerConfigMiddleware struct {
	md *authMiddleware
}

func (w *webServerConfigMiddleware) GetOptions() httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.BifrostGroupMeta]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerConfigs().GetOptions())
}

func (w *webServerConfigMiddleware) GetConfigTextLines() httpclientv1.ClientBuilder[v1.WebServerOptions, modelclientv1.ResponseBody[string]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerConfigs().GetConfigTextLines())
}

func (w *webServerConfigMiddleware) GetContextTextLines() httpclientv1.ClientBuilder[v1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[string]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerConfigs().GetContextTextLines())
}

func (w *webServerConfigMiddleware) GetConfig() httpclientv1.ClientBuilder[v1.WebServerOptions, modelclientv1.ResponseBody[*modelclientv1.WebServerConfig]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerConfigs().GetConfig())
}

func (w *webServerConfigMiddleware) GetIncludedConfigs() httpclientv1.ClientBuilder[v1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[[]string]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerConfigs().GetIncludedConfigs())
}

func (w *webServerConfigMiddleware) SearchContextPositions() httpclientv1.ClientBuilder[v1.WebServerConfigContextPosSearchOptions, modelclientv1.ResponseBody[[]v1.ConfigContextPos]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerConfigs().SearchContextPositions())
}

func (w *webServerConfigMiddleware) InsertWithClone() httpclientv1.ClientBuilder[v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerConfigs().InsertWithClone())
}

func (w *webServerConfigMiddleware) InsertWithNew() httpclientv1.ClientBuilder[v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerConfigs().InsertWithNew())
}

func (w *webServerConfigMiddleware) Remove() httpclientv1.ClientBuilder[v1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerConfigs().Remove())
}

func (w *webServerConfigMiddleware) ModifyContextValue() httpclientv1.ClientBuilder[v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerConfigs().ModifyContextValue())
}

func (w *webServerConfigMiddleware) ModifyWithClone() httpclientv1.ClientBuilder[v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerConfigs().ModifyWithClone())
}

func (w *webServerConfigMiddleware) ChangeContextEnabledState() httpclientv1.ClientBuilder[v1.WebServerConfigContextUpdateOptions[v1.ConfigContextEnabledStateMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerConfigs().ChangeContextEnabledState())
}

func (w *webServerConfigMiddleware) ModifyWithNew() httpclientv1.ClientBuilder[v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerConfigs().ModifyWithNew())
}

func (w *webServerConfigMiddleware) Move() httpclientv1.ClientBuilder[v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerConfigs().Move())
}

func (w *webServerConfigMiddleware) UpdateConfig() httpclientv1.ClientBuilder[*v1.WebServerConfigUpdateOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	return wrapWithAuth(w.md, w.md.txp.WebServerConfigs().UpdateConfig())
}

func newWebServerConfigMiddleware(a *authMiddleware) *webServerConfigMiddleware {
	return &webServerConfigMiddleware{md: a}
}
