package endpoint

import (
	v1 "github.com/tanganyu1114/heimdallr-reborn/api/heimdallr_api/v1"
	"testing"

	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
	"github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/transport"

	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_newWebServerConfigEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockFactory(ctrl)
	mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockTransport.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport).AnyTimes()

	type args struct {
		f *factory
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "creates web server config endpoints",
			args: args{
				f: &factory{
					transport: mockTransport,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newWebServerConfigEndpoints(tt.args.f)
			if got == nil {
				t.Errorf("newWebServerConfigEndpoints() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigEndpoints_ChangeContextEnabledState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().ChangeContextEnabledState().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerConfigTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], httpclientv1.NilBody]
	}{
		{
			name:   "returns change context enabled state endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigEndpoints{
				transport: tt.fields.transport,
			}
			got := w.ChangeContextEnabledState()
			if got == nil {
				t.Errorf("ChangeContextEnabledState() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigEndpoints_GetConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[*modelclientv1.WebServerConfig]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.WebServerOptions, modelclientv1.ResponseBody[*modelclientv1.WebServerConfig]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[*modelclientv1.WebServerConfig]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().GetConfig().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerConfigTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[*modelclientv1.WebServerConfig]]
	}{
		{
			name:   "returns get config endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigEndpoints{
				transport: tt.fields.transport,
			}
			got := w.GetConfig()
			if got == nil {
				t.Errorf("GetConfig() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigEndpoints_GetConfigTextLines(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[string]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.WebServerOptions, modelclientv1.ResponseBody[string]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[string]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().GetConfigTextLines().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerConfigTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[string]]
	}{
		{
			name:   "returns get config text lines endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigEndpoints{
				transport: tt.fields.transport,
			}
			got := w.GetConfigTextLines()
			if got == nil {
				t.Errorf("GetConfigTextLines() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigEndpoints_GetContextTextLines(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[string]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[string]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[string]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().GetContextTextLines().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerConfigTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[string]]
	}{
		{
			name:   "returns get context text lines endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigEndpoints{
				transport: tt.fields.transport,
			}
			got := w.GetContextTextLines()
			if got == nil {
				t.Errorf("GetContextTextLines() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigEndpoints_GetIncludedConfigs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[[]string]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[[]string]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[[]string]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().GetIncludedConfigs().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerConfigTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[[]string]]
	}{
		{
			name:   "returns get included configs endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigEndpoints{
				transport: tt.fields.transport,
			}
			got := w.GetIncludedConfigs()
			if got == nil {
				t.Errorf("GetIncludedConfigs() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigEndpoints_GetOptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.BifrostGroupMeta]](ctrl)
	mockClient := httpclientv1.NewMockClient[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.BifrostGroupMeta]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.BifrostGroupMeta]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().GetOptions().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerConfigTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.BifrostGroupMeta]]
	}{
		{
			name:   "returns get options endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigEndpoints{
				transport: tt.fields.transport,
			}
			got := w.GetOptions()
			if got == nil {
				t.Errorf("GetOptions() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigEndpoints_InsertWithClone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().InsertWithClone().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerConfigTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}{
		{
			name:   "returns insert with clone endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigEndpoints{
				transport: tt.fields.transport,
			}
			got := w.InsertWithClone()
			if got == nil {
				t.Errorf("InsertWithClone() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigEndpoints_InsertWithNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().InsertWithNew().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerConfigTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}{
		{
			name:   "returns insert with new endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigEndpoints{
				transport: tt.fields.transport,
			}
			got := w.InsertWithNew()
			if got == nil {
				t.Errorf("InsertWithNew() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigEndpoints_ModifyContextValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().ModifyContextValue().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerConfigTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}{
		{
			name:   "returns modify context value endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigEndpoints{
				transport: tt.fields.transport,
			}
			got := w.ModifyContextValue()
			if got == nil {
				t.Errorf("ModifyContextValue() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigEndpoints_ModifyWithClone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().ModifyWithClone().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerConfigTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}{
		{
			name:   "returns modify with clone endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigEndpoints{
				transport: tt.fields.transport,
			}
			got := w.ModifyWithClone()
			if got == nil {
				t.Errorf("ModifyWithClone() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigEndpoints_ModifyWithNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().ModifyWithNew().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerConfigTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}{
		{
			name:   "returns modify with new endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigEndpoints{
				transport: tt.fields.transport,
			}
			got := w.ModifyWithNew()
			if got == nil {
				t.Errorf("ModifyWithNew() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigEndpoints_Move(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().Move().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerConfigTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}{
		{
			name:   "returns move endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigEndpoints{
				transport: tt.fields.transport,
			}
			got := w.Move()
			if got == nil {
				t.Errorf("Move() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigEndpoints_Remove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().Remove().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerConfigTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}{
		{
			name:   "returns remove endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigEndpoints{
				transport: tt.fields.transport,
			}
			got := w.Remove()
			if got == nil {
				t.Errorf("Remove() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigEndpoints_SearchContextPositions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextPosSearchOptions, modelclientv1.ResponseBody[[]metav1.ConfigContextPos]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.WebServerConfigContextPosSearchOptions, modelclientv1.ResponseBody[[]metav1.ConfigContextPos]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextPosSearchOptions, modelclientv1.ResponseBody[[]metav1.ConfigContextPos]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().SearchContextPositions().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerConfigTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.WebServerConfigContextPosSearchOptions, modelclientv1.ResponseBody[[]metav1.ConfigContextPos]]
	}{
		{
			name:   "returns search context positions endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigEndpoints{
				transport: tt.fields.transport,
			}
			got := w.SearchContextPositions()
			if got == nil {
				t.Errorf("SearchContextPositions() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerConfigEndpoints_UpdateConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[*metav1.WebServerConfigUpdateOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockClient := httpclientv1.NewMockClient[*metav1.WebServerConfigUpdateOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[*metav1.WebServerConfigUpdateOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().UpdateConfig().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerConfigTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[*metav1.WebServerConfigUpdateOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}{
		{
			name:   "returns update config endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigEndpoints{
				transport: tt.fields.transport,
			}
			got := w.UpdateConfig()
			if got == nil {
				t.Errorf("UpdateConfig() = nil, want non-nil")
			}
		})
	}
}
