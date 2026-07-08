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

func Test_webServerStatisticsMiddleware_ConnectivityCheckOfProxyService(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "connectivity check of proxy service middleware",
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
			mockWebServerStatisticsTransport := transport.NewMockWebServerStatisticsTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]](ctrl)

			mockFactory.EXPECT().WebServerStatistics().Return(mockWebServerStatisticsTransport)
			mockWebServerStatisticsTransport.EXPECT().ConnectivityCheckOfProxyService().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerStatisticsMiddleware{
				md: authMw,
			}
			got := w.ConnectivityCheckOfProxyService()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerStatisticsMiddleware_ExportProxyServiceInfoToExcel(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "export proxy service info to excel middleware",
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
			mockWebServerStatisticsTransport := transport.NewMockWebServerStatisticsTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]byte]](ctrl)

			mockFactory.EXPECT().WebServerStatistics().Return(mockWebServerStatisticsTransport)
			mockWebServerStatisticsTransport.EXPECT().ExportProxyServiceInfoToExcel().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerStatisticsMiddleware{
				md: authMw,
			}
			got := w.ExportProxyServiceInfoToExcel()
			assert.NotNil(t, got)
		})
	}
}

func Test_webServerStatisticsMiddleware_GetProxyServiceInfo(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "get proxy service info middleware",
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
			mockWebServerStatisticsTransport := transport.NewMockWebServerStatisticsTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]](ctrl)

			mockFactory.EXPECT().WebServerStatistics().Return(mockWebServerStatisticsTransport)
			mockWebServerStatisticsTransport.EXPECT().GetProxyServiceInfo().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerStatisticsMiddleware{
				md: authMw,
			}
			got := w.GetProxyServiceInfo()
			assert.NotNil(t, got)
		})
	}
}
