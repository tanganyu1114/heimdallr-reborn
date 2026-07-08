package auth

import (
	"testing"

	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	modelclientv1 "gin-vue-admin/pkg/client/v1/model"
	"gin-vue-admin/pkg/client/v1/transport"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_webServerConfigMiddleware_ChangeContextEnabledState(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "change context enabled state middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], httpclientv1.NilBody](ctrl)

			mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport)
			mockWebServerConfigTransport.EXPECT().ChangeContextEnabledState().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerConfigMiddleware{
				md: authMw,
			}
			got := w.ChangeContextEnabledState()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerConfigMiddleware_GetConfig(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "get config middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[*modelclientv1.WebServerConfig]](ctrl)

			mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport)
			mockWebServerConfigTransport.EXPECT().GetConfig().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerConfigMiddleware{
				md: authMw,
			}
			got := w.GetConfig()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerConfigMiddleware_GetConfigTextLines(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "get config text lines middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[string]](ctrl)

			mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport)
			mockWebServerConfigTransport.EXPECT().GetConfigTextLines().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerConfigMiddleware{
				md: authMw,
			}
			got := w.GetConfigTextLines()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerConfigMiddleware_GetContextTextLines(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "get context text lines middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[string]](ctrl)

			mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport)
			mockWebServerConfigTransport.EXPECT().GetContextTextLines().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerConfigMiddleware{
				md: authMw,
			}
			got := w.GetContextTextLines()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerConfigMiddleware_GetIncludedConfigs(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "get included configs middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[[]string]](ctrl)

			mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport)
			mockWebServerConfigTransport.EXPECT().GetIncludedConfigs().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerConfigMiddleware{
				md: authMw,
			}
			got := w.GetIncludedConfigs()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerConfigMiddleware_GetOptions(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "get options middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.BifrostGroupMeta]](ctrl)

			mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport)
			mockWebServerConfigTransport.EXPECT().GetOptions().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerConfigMiddleware{
				md: authMw,
			}
			got := w.GetOptions()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerConfigMiddleware_InsertWithClone(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "insert with clone middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody](ctrl)

			mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport)
			mockWebServerConfigTransport.EXPECT().InsertWithClone().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerConfigMiddleware{
				md: authMw,
			}
			got := w.InsertWithClone()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerConfigMiddleware_InsertWithNew(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "insert with new middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody](ctrl)

			mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport)
			mockWebServerConfigTransport.EXPECT().InsertWithNew().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerConfigMiddleware{
				md: authMw,
			}
			got := w.InsertWithNew()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerConfigMiddleware_ModifyContextValue(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "modify context value middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody](ctrl)

			mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport)
			mockWebServerConfigTransport.EXPECT().ModifyContextValue().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerConfigMiddleware{
				md: authMw,
			}
			got := w.ModifyContextValue()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerConfigMiddleware_ModifyWithClone(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "modify with clone middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody](ctrl)

			mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport)
			mockWebServerConfigTransport.EXPECT().ModifyWithClone().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerConfigMiddleware{
				md: authMw,
			}
			got := w.ModifyWithClone()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerConfigMiddleware_ModifyWithNew(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "modify with new middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody](ctrl)

			mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport)
			mockWebServerConfigTransport.EXPECT().ModifyWithNew().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerConfigMiddleware{
				md: authMw,
			}
			got := w.ModifyWithNew()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerConfigMiddleware_Move(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "move middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody](ctrl)

			mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport)
			mockWebServerConfigTransport.EXPECT().Move().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerConfigMiddleware{
				md: authMw,
			}
			got := w.Move()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerConfigMiddleware_Remove(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "remove middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigTargetContextOptions, httpclientv1.NilBody](ctrl)

			mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport)
			mockWebServerConfigTransport.EXPECT().Remove().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerConfigMiddleware{
				md: authMw,
			}
			got := w.Remove()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerConfigMiddleware_SearchContextPositions(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "search context positions middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerConfigContextPosSearchOptions, modelclientv1.ResponseBody[[]metav1.ConfigContextPos]](ctrl)

			mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport)
			mockWebServerConfigTransport.EXPECT().SearchContextPositions().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerConfigMiddleware{
				md: authMw,
			}
			got := w.SearchContextPositions()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerConfigMiddleware_UpdateConfig(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "update config middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[*metav1.WebServerConfigUpdateOptions, httpclientv1.NilBody](ctrl)

			mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport)
			mockWebServerConfigTransport.EXPECT().UpdateConfig().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerConfigMiddleware{
				md: authMw,
			}
			got := w.UpdateConfig()
			assert.NotNil(t, got)
		})
	}
}
