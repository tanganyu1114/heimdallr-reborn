package endpoint

import (
	"testing"

	"github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/transport"

	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_newWebServerStatisticsEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockFactory(ctrl)
	mockWebServerStatisticsTransport := transport.NewMockWebServerStatisticsTransport(ctrl)
	mockTransport.EXPECT().WebServerStatistics().Return(mockWebServerStatisticsTransport).AnyTimes()

	type args struct {
		f *factory
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "creates web server statistics endpoints",
			args: args{
				f: &factory{
					transport: mockTransport,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newWebServerStatisticsEndpoints(tt.args.f)
			if got == nil {
				t.Errorf("newWebServerStatisticsEndpoints() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerStatisticsEndpoints_ConnectivityCheckOfProxyService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerStatisticsTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[v1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]](ctrl)
	mockClient := httpclientv1.NewMockClient[v1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[v1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().ConnectivityCheckOfProxyService().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerStatisticsTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[v1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]]
	}{
		{
			name:   "returns connectivity check of proxy service endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerStatisticsEndpoints{
				transport: tt.fields.transport,
			}
			got := w.ConnectivityCheckOfProxyService()
			if got == nil {
				t.Errorf("ConnectivityCheckOfProxyService() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerStatisticsEndpoints_ExportProxyServiceInfoToExcel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerStatisticsTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[v1.WebServerOptions, modelclientv1.ResponseBody[[]byte]](ctrl)
	mockClient := httpclientv1.NewMockClient[v1.WebServerOptions, modelclientv1.ResponseBody[[]byte]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[v1.WebServerOptions, modelclientv1.ResponseBody[[]byte]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().ExportProxyServiceInfoToExcel().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerStatisticsTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[v1.WebServerOptions, modelclientv1.ResponseBody[[]byte]]
	}{
		{
			name:   "returns export proxy service info to excel endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerStatisticsEndpoints{
				transport: tt.fields.transport,
			}
			got := w.ExportProxyServiceInfoToExcel()
			if got == nil {
				t.Errorf("ExportProxyServiceInfoToExcel() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerStatisticsEndpoints_GetProxyServiceInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerStatisticsTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[v1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]](ctrl)
	mockClient := httpclientv1.NewMockClient[v1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[v1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().GetProxyServiceInfo().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerStatisticsTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[v1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]]
	}{
		{
			name:   "returns get proxy service info endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerStatisticsEndpoints{
				transport: tt.fields.transport,
			}
			got := w.GetProxyServiceInfo()
			if got == nil {
				t.Errorf("GetProxyServiceInfo() = nil, want non-nil")
			}
		})
	}
}
