package auth

import (
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

type webServerConfigMiddleware struct {
	md *authMiddleware
}

func (w *webServerConfigMiddleware) GetOptions() httpclientv1.ClientBuilder[httpclientv1.NilBody, []v1.BifrostGroupMeta] {
	return applyAuthOptions(w.md, w.md.txp.WebServerConfigs().GetOptions())
}

func (w *webServerConfigMiddleware) GetConfigTextLines() httpclientv1.ClientBuilder[metav1.WebServerOptions, string] {
	return applyAuthOptions(w.md, w.md.txp.WebServerConfigs().GetConfigTextLines())
}

func (w *webServerConfigMiddleware) GetContextTextLines() httpclientv1.ClientBuilder[metav1.WebServerConfigTargetContextOptions, string] {
	return applyAuthOptions(w.md, w.md.txp.WebServerConfigs().GetContextTextLines())
}

func (w *webServerConfigMiddleware) GetConfig() httpclientv1.ClientBuilder[metav1.WebServerOptions, *metav1.WebServerConfig] {
	return applyAuthOptions(w.md, w.md.txp.WebServerConfigs().GetConfig())
}

func (w *webServerConfigMiddleware) GetIncludedConfigs() httpclientv1.ClientBuilder[metav1.WebServerConfigTargetContextOptions, []string] {
	return applyAuthOptions(w.md, w.md.txp.WebServerConfigs().GetIncludedConfigs())
}

func (w *webServerConfigMiddleware) SearchContextPositions() httpclientv1.ClientBuilder[metav1.WebServerConfigContextPosSearchOptions, []metav1.ConfigContextPos] {
	return applyAuthOptions(w.md, w.md.txp.WebServerConfigs().SearchContextPositions())
}

func (w *webServerConfigMiddleware) InsertWithClone() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody] {
	return applyAuthOptions(w.md, w.md.txp.WebServerConfigs().InsertWithClone())
}

func (w *webServerConfigMiddleware) InsertWithNew() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody] {
	return applyAuthOptions(w.md, w.md.txp.WebServerConfigs().InsertWithNew())
}

func (w *webServerConfigMiddleware) Remove() httpclientv1.ClientBuilder[metav1.WebServerConfigTargetContextOptions, httpclientv1.NilBody] {
	return applyAuthOptions(w.md, w.md.txp.WebServerConfigs().Remove())
}

func (w *webServerConfigMiddleware) ModifyContextValue() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody] {
	return applyAuthOptions(w.md, w.md.txp.WebServerConfigs().ModifyContextValue())
}

func (w *webServerConfigMiddleware) ModifyWithClone() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody] {
	return applyAuthOptions(w.md, w.md.txp.WebServerConfigs().ModifyWithClone())
}

func (w *webServerConfigMiddleware) ChangeContextEnabledState() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], httpclientv1.NilBody] {
	return applyAuthOptions(w.md, w.md.txp.WebServerConfigs().ChangeContextEnabledState())
}

func (w *webServerConfigMiddleware) ModifyWithNew() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody] {
	return applyAuthOptions(w.md, w.md.txp.WebServerConfigs().ModifyWithNew())
}

func (w *webServerConfigMiddleware) Move() httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody] {
	return applyAuthOptions(w.md, w.md.txp.WebServerConfigs().Move())
}

func newWebServerConfigMiddleware(a *authMiddleware) *webServerConfigMiddleware {
	return &webServerConfigMiddleware{md: a}
}
