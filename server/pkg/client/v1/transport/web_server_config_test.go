package transport

import (
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"reflect"
	"testing"

	metav1 "gin-vue-admin/internal/pkg/meta/v1"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_newWebServerConfigTransport(t *testing.T) {
	type args struct {
		transport *transport
	}
	tests := []struct {
		name string
		args args
		want WebServerConfigTransport
	}{
		{
			name: "creates web server config transport with valid base URL",
			args: args{
				transport: &transport{
					baseURL: "http://localhost:8080",
				},
			},
		},
		{
			name: "creates web server config transport with empty base URL",
			args: args{
				transport: &transport{
					baseURL: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newWebServerConfigTransport(tt.args.transport)
			if got == nil {
				t.Errorf("newWebServerConfigTransport() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigTransport_ChangeContextEnabledState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], httpclientv1.NilBody](ctrl)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], httpclientv1.NilBody]
	}{
		{
			name: "returns change context enabled state client",
			fields: fields{
				changeContextEnabledStateClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigTransport{
				getOptionsClient:                tt.fields.getOptionsClient,
				getConfigTextLinesClient:        tt.fields.getConfigTextLinesClient,
				getContextTextLinesClient:       tt.fields.getContextTextLinesClient,
				getConfigClient:                 tt.fields.getConfigClient,
				getIncludedConfigsClient:        tt.fields.getIncludedConfigsClient,
				searchContextPositionsClient:    tt.fields.searchContextPositionsClient,
				insertWithCloneClient:           tt.fields.insertWithCloneClient,
				insertWithNewClient:             tt.fields.insertWithNewClient,
				removeClient:                    tt.fields.removeClient,
				modifyContextValueClient:        tt.fields.modifyContextValueClient,
				modifyWithCloneClient:           tt.fields.modifyWithCloneClient,
				changeContextEnabledStateClient: tt.fields.changeContextEnabledStateClient,
				modifyWithNewClient:             tt.fields.modifyWithNewClient,
				moveClient:                      tt.fields.moveClient,
			}
			if got := w.ChangeContextEnabledState(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChangeContextEnabledState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigTransport_GetConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerOptions, *metav1.WebServerConfig](ctrl)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerOptions, *metav1.WebServerConfig]
	}{
		{
			name: "returns get config client",
			fields: fields{
				getConfigClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigTransport{
				getOptionsClient:                tt.fields.getOptionsClient,
				getConfigTextLinesClient:        tt.fields.getConfigTextLinesClient,
				getContextTextLinesClient:       tt.fields.getContextTextLinesClient,
				getConfigClient:                 tt.fields.getConfigClient,
				getIncludedConfigsClient:        tt.fields.getIncludedConfigsClient,
				searchContextPositionsClient:    tt.fields.searchContextPositionsClient,
				insertWithCloneClient:           tt.fields.insertWithCloneClient,
				insertWithNewClient:             tt.fields.insertWithNewClient,
				removeClient:                    tt.fields.removeClient,
				modifyContextValueClient:        tt.fields.modifyContextValueClient,
				modifyWithCloneClient:           tt.fields.modifyWithCloneClient,
				changeContextEnabledStateClient: tt.fields.changeContextEnabledStateClient,
				modifyWithNewClient:             tt.fields.modifyWithNewClient,
				moveClient:                      tt.fields.moveClient,
			}
			if got := w.GetConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigTransport_GetConfigTextLines(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerOptions, string](ctrl)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerOptions, string]
	}{
		{
			name: "returns get config text lines client",
			fields: fields{
				getConfigTextLinesClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigTransport{
				getOptionsClient:                tt.fields.getOptionsClient,
				getConfigTextLinesClient:        tt.fields.getConfigTextLinesClient,
				getContextTextLinesClient:       tt.fields.getContextTextLinesClient,
				getConfigClient:                 tt.fields.getConfigClient,
				getIncludedConfigsClient:        tt.fields.getIncludedConfigsClient,
				searchContextPositionsClient:    tt.fields.searchContextPositionsClient,
				insertWithCloneClient:           tt.fields.insertWithCloneClient,
				insertWithNewClient:             tt.fields.insertWithNewClient,
				removeClient:                    tt.fields.removeClient,
				modifyContextValueClient:        tt.fields.modifyContextValueClient,
				modifyWithCloneClient:           tt.fields.modifyWithCloneClient,
				changeContextEnabledStateClient: tt.fields.changeContextEnabledStateClient,
				modifyWithNewClient:             tt.fields.modifyWithNewClient,
				moveClient:                      tt.fields.moveClient,
			}
			if got := w.GetConfigTextLines(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfigTextLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigTransport_GetContextTextLines(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigTargetContextOptions, string](ctrl)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerConfigTargetContextOptions, string]
	}{
		{
			name: "returns get context text lines client",
			fields: fields{
				getContextTextLinesClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigTransport{
				getOptionsClient:                tt.fields.getOptionsClient,
				getConfigTextLinesClient:        tt.fields.getConfigTextLinesClient,
				getContextTextLinesClient:       tt.fields.getContextTextLinesClient,
				getConfigClient:                 tt.fields.getConfigClient,
				getIncludedConfigsClient:        tt.fields.getIncludedConfigsClient,
				searchContextPositionsClient:    tt.fields.searchContextPositionsClient,
				insertWithCloneClient:           tt.fields.insertWithCloneClient,
				insertWithNewClient:             tt.fields.insertWithNewClient,
				removeClient:                    tt.fields.removeClient,
				modifyContextValueClient:        tt.fields.modifyContextValueClient,
				modifyWithCloneClient:           tt.fields.modifyWithCloneClient,
				changeContextEnabledStateClient: tt.fields.changeContextEnabledStateClient,
				modifyWithNewClient:             tt.fields.modifyWithNewClient,
				moveClient:                      tt.fields.moveClient,
			}
			if got := w.GetContextTextLines(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContextTextLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigTransport_GetIncludedConfigs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigTargetContextOptions, []string](ctrl)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerConfigTargetContextOptions, []string]
	}{
		{
			name: "returns get included configs client",
			fields: fields{
				getIncludedConfigsClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigTransport{
				getOptionsClient:                tt.fields.getOptionsClient,
				getConfigTextLinesClient:        tt.fields.getConfigTextLinesClient,
				getContextTextLinesClient:       tt.fields.getContextTextLinesClient,
				getConfigClient:                 tt.fields.getConfigClient,
				getIncludedConfigsClient:        tt.fields.getIncludedConfigsClient,
				searchContextPositionsClient:    tt.fields.searchContextPositionsClient,
				insertWithCloneClient:           tt.fields.insertWithCloneClient,
				insertWithNewClient:             tt.fields.insertWithNewClient,
				removeClient:                    tt.fields.removeClient,
				modifyContextValueClient:        tt.fields.modifyContextValueClient,
				modifyWithCloneClient:           tt.fields.modifyWithCloneClient,
				changeContextEnabledStateClient: tt.fields.changeContextEnabledStateClient,
				modifyWithNewClient:             tt.fields.modifyWithNewClient,
				moveClient:                      tt.fields.moveClient,
			}
			if got := w.GetIncludedConfigs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIncludedConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigTransport_GetOptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[httpclientv1.NilBody, []v1.BifrostGroupMeta](ctrl)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[httpclientv1.NilBody, []v1.BifrostGroupMeta]
	}{
		{
			name: "returns get options client",
			fields: fields{
				getOptionsClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigTransport{
				getOptionsClient:                tt.fields.getOptionsClient,
				getConfigTextLinesClient:        tt.fields.getConfigTextLinesClient,
				getContextTextLinesClient:       tt.fields.getContextTextLinesClient,
				getConfigClient:                 tt.fields.getConfigClient,
				getIncludedConfigsClient:        tt.fields.getIncludedConfigsClient,
				searchContextPositionsClient:    tt.fields.searchContextPositionsClient,
				insertWithCloneClient:           tt.fields.insertWithCloneClient,
				insertWithNewClient:             tt.fields.insertWithNewClient,
				removeClient:                    tt.fields.removeClient,
				modifyContextValueClient:        tt.fields.modifyContextValueClient,
				modifyWithCloneClient:           tt.fields.modifyWithCloneClient,
				changeContextEnabledStateClient: tt.fields.changeContextEnabledStateClient,
				modifyWithNewClient:             tt.fields.modifyWithNewClient,
				moveClient:                      tt.fields.moveClient,
			}
			if got := w.GetOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigTransport_InsertWithClone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody](ctrl)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody]
	}{
		{
			name: "returns insert with clone client",
			fields: fields{
				insertWithCloneClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigTransport{
				getOptionsClient:                tt.fields.getOptionsClient,
				getConfigTextLinesClient:        tt.fields.getConfigTextLinesClient,
				getContextTextLinesClient:       tt.fields.getContextTextLinesClient,
				getConfigClient:                 tt.fields.getConfigClient,
				getIncludedConfigsClient:        tt.fields.getIncludedConfigsClient,
				searchContextPositionsClient:    tt.fields.searchContextPositionsClient,
				insertWithCloneClient:           tt.fields.insertWithCloneClient,
				insertWithNewClient:             tt.fields.insertWithNewClient,
				removeClient:                    tt.fields.removeClient,
				modifyContextValueClient:        tt.fields.modifyContextValueClient,
				modifyWithCloneClient:           tt.fields.modifyWithCloneClient,
				changeContextEnabledStateClient: tt.fields.changeContextEnabledStateClient,
				modifyWithNewClient:             tt.fields.modifyWithNewClient,
				moveClient:                      tt.fields.moveClient,
			}
			if got := w.InsertWithClone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertWithClone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigTransport_InsertWithNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody](ctrl)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody]
	}{
		{
			name: "returns insert with new client",
			fields: fields{
				insertWithNewClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigTransport{
				getOptionsClient:                tt.fields.getOptionsClient,
				getConfigTextLinesClient:        tt.fields.getConfigTextLinesClient,
				getContextTextLinesClient:       tt.fields.getContextTextLinesClient,
				getConfigClient:                 tt.fields.getConfigClient,
				getIncludedConfigsClient:        tt.fields.getIncludedConfigsClient,
				searchContextPositionsClient:    tt.fields.searchContextPositionsClient,
				insertWithCloneClient:           tt.fields.insertWithCloneClient,
				insertWithNewClient:             tt.fields.insertWithNewClient,
				removeClient:                    tt.fields.removeClient,
				modifyContextValueClient:        tt.fields.modifyContextValueClient,
				modifyWithCloneClient:           tt.fields.modifyWithCloneClient,
				changeContextEnabledStateClient: tt.fields.changeContextEnabledStateClient,
				modifyWithNewClient:             tt.fields.modifyWithNewClient,
				moveClient:                      tt.fields.moveClient,
			}
			if got := w.InsertWithNew(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertWithNew() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigTransport_ModifyContextValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody](ctrl)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody]
	}{
		{
			name: "returns modify context value client",
			fields: fields{
				modifyContextValueClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigTransport{
				getOptionsClient:                tt.fields.getOptionsClient,
				getConfigTextLinesClient:        tt.fields.getConfigTextLinesClient,
				getContextTextLinesClient:       tt.fields.getContextTextLinesClient,
				getConfigClient:                 tt.fields.getConfigClient,
				getIncludedConfigsClient:        tt.fields.getIncludedConfigsClient,
				searchContextPositionsClient:    tt.fields.searchContextPositionsClient,
				insertWithCloneClient:           tt.fields.insertWithCloneClient,
				insertWithNewClient:             tt.fields.insertWithNewClient,
				removeClient:                    tt.fields.removeClient,
				modifyContextValueClient:        tt.fields.modifyContextValueClient,
				modifyWithCloneClient:           tt.fields.modifyWithCloneClient,
				changeContextEnabledStateClient: tt.fields.changeContextEnabledStateClient,
				modifyWithNewClient:             tt.fields.modifyWithNewClient,
				moveClient:                      tt.fields.moveClient,
			}
			if got := w.ModifyContextValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModifyContextValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigTransport_ModifyWithClone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody](ctrl)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody]
	}{
		{
			name: "returns modify with clone client",
			fields: fields{
				modifyWithCloneClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigTransport{
				getOptionsClient:                tt.fields.getOptionsClient,
				getConfigTextLinesClient:        tt.fields.getConfigTextLinesClient,
				getContextTextLinesClient:       tt.fields.getContextTextLinesClient,
				getConfigClient:                 tt.fields.getConfigClient,
				getIncludedConfigsClient:        tt.fields.getIncludedConfigsClient,
				searchContextPositionsClient:    tt.fields.searchContextPositionsClient,
				insertWithCloneClient:           tt.fields.insertWithCloneClient,
				insertWithNewClient:             tt.fields.insertWithNewClient,
				removeClient:                    tt.fields.removeClient,
				modifyContextValueClient:        tt.fields.modifyContextValueClient,
				modifyWithCloneClient:           tt.fields.modifyWithCloneClient,
				changeContextEnabledStateClient: tt.fields.changeContextEnabledStateClient,
				modifyWithNewClient:             tt.fields.modifyWithNewClient,
				moveClient:                      tt.fields.moveClient,
			}
			if got := w.ModifyWithClone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModifyWithClone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigTransport_ModifyWithNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody](ctrl)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody]
	}{
		{
			name: "returns modify with new client",
			fields: fields{
				modifyWithNewClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigTransport{
				getOptionsClient:                tt.fields.getOptionsClient,
				getConfigTextLinesClient:        tt.fields.getConfigTextLinesClient,
				getContextTextLinesClient:       tt.fields.getContextTextLinesClient,
				getConfigClient:                 tt.fields.getConfigClient,
				getIncludedConfigsClient:        tt.fields.getIncludedConfigsClient,
				searchContextPositionsClient:    tt.fields.searchContextPositionsClient,
				insertWithCloneClient:           tt.fields.insertWithCloneClient,
				insertWithNewClient:             tt.fields.insertWithNewClient,
				removeClient:                    tt.fields.removeClient,
				modifyContextValueClient:        tt.fields.modifyContextValueClient,
				modifyWithCloneClient:           tt.fields.modifyWithCloneClient,
				changeContextEnabledStateClient: tt.fields.changeContextEnabledStateClient,
				modifyWithNewClient:             tt.fields.modifyWithNewClient,
				moveClient:                      tt.fields.moveClient,
			}
			if got := w.ModifyWithNew(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModifyWithNew() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigTransport_Move(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody](ctrl)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody]
	}{
		{
			name: "returns move client",
			fields: fields{
				moveClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigTransport{
				getOptionsClient:                tt.fields.getOptionsClient,
				getConfigTextLinesClient:        tt.fields.getConfigTextLinesClient,
				getContextTextLinesClient:       tt.fields.getContextTextLinesClient,
				getConfigClient:                 tt.fields.getConfigClient,
				getIncludedConfigsClient:        tt.fields.getIncludedConfigsClient,
				searchContextPositionsClient:    tt.fields.searchContextPositionsClient,
				insertWithCloneClient:           tt.fields.insertWithCloneClient,
				insertWithNewClient:             tt.fields.insertWithNewClient,
				removeClient:                    tt.fields.removeClient,
				modifyContextValueClient:        tt.fields.modifyContextValueClient,
				modifyWithCloneClient:           tt.fields.modifyWithCloneClient,
				changeContextEnabledStateClient: tt.fields.changeContextEnabledStateClient,
				modifyWithNewClient:             tt.fields.modifyWithNewClient,
				moveClient:                      tt.fields.moveClient,
			}
			if got := w.Move(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Move() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigTransport_Remove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigTargetContextOptions, httpclientv1.NilBody](ctrl)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerConfigTargetContextOptions, httpclientv1.NilBody]
	}{
		{
			name: "returns remove client",
			fields: fields{
				removeClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigTransport{
				getOptionsClient:                tt.fields.getOptionsClient,
				getConfigTextLinesClient:        tt.fields.getConfigTextLinesClient,
				getContextTextLinesClient:       tt.fields.getContextTextLinesClient,
				getConfigClient:                 tt.fields.getConfigClient,
				getIncludedConfigsClient:        tt.fields.getIncludedConfigsClient,
				searchContextPositionsClient:    tt.fields.searchContextPositionsClient,
				insertWithCloneClient:           tt.fields.insertWithCloneClient,
				insertWithNewClient:             tt.fields.insertWithNewClient,
				removeClient:                    tt.fields.removeClient,
				modifyContextValueClient:        tt.fields.modifyContextValueClient,
				modifyWithCloneClient:           tt.fields.modifyWithCloneClient,
				changeContextEnabledStateClient: tt.fields.changeContextEnabledStateClient,
				modifyWithNewClient:             tt.fields.modifyWithNewClient,
				moveClient:                      tt.fields.moveClient,
			}
			if got := w.Remove(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigTransport_SearchContextPositions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextPosSearchOptions, []metav1.ConfigContextPos](ctrl)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerConfigContextPosSearchOptions, []metav1.ConfigContextPos]
	}{
		{
			name: "returns search context positions client",
			fields: fields{
				searchContextPositionsClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigTransport{
				getOptionsClient:                tt.fields.getOptionsClient,
				getConfigTextLinesClient:        tt.fields.getConfigTextLinesClient,
				getContextTextLinesClient:       tt.fields.getContextTextLinesClient,
				getConfigClient:                 tt.fields.getConfigClient,
				getIncludedConfigsClient:        tt.fields.getIncludedConfigsClient,
				searchContextPositionsClient:    tt.fields.searchContextPositionsClient,
				insertWithCloneClient:           tt.fields.insertWithCloneClient,
				insertWithNewClient:             tt.fields.insertWithNewClient,
				removeClient:                    tt.fields.removeClient,
				modifyContextValueClient:        tt.fields.modifyContextValueClient,
				modifyWithCloneClient:           tt.fields.modifyWithCloneClient,
				changeContextEnabledStateClient: tt.fields.changeContextEnabledStateClient,
				modifyWithNewClient:             tt.fields.modifyWithNewClient,
				moveClient:                      tt.fields.moveClient,
			}
			if got := w.SearchContextPositions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchContextPositions() = %v, want %v", got, tt.want)
			}
		})
	}
}
