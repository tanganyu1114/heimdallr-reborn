package transport

import (
	"reflect"
	"testing"

	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	modelclientv1 "gin-vue-admin/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_newWebServerStatisticsTransport(t *testing.T) {
	type args struct {
		transport *transport
	}
	tests := []struct {
		name string
		args args
		want WebServerStatisticsTransport
	}{
		{
			name: "creates web server statistics transport with valid base URL",
			args: args{
				transport: &transport{
					baseURL: "http://localhost:8080",
				},
			},
		},
		{
			name: "creates web server statistics transport with empty base URL",
			args: args{
				transport: &transport{
					baseURL: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newWebServerStatisticsTransport(tt.args.transport)
			if got == nil {
				t.Errorf("newWebServerStatisticsTransport() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerStatisticsTransport_ConnectivityCheckOfProxyService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]](ctrl)

	type fields struct {
		getProxyServiceInfoClient             httpclientv1.ClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]]
		connectivityCheckOfProxyServiceClient httpclientv1.ClientBuilder[metav1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]]
		exportProxyServiceInfoToExcelClient   httpclientv1.ClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]byte]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]]
	}{
		{
			name: "returns connectivity check of proxy service client",
			fields: fields{
				connectivityCheckOfProxyServiceClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerStatisticsTransport{
				getProxyServiceInfoClient:             tt.fields.getProxyServiceInfoClient,
				connectivityCheckOfProxyServiceClient: tt.fields.connectivityCheckOfProxyServiceClient,
				exportProxyServiceInfoToExcelClient:   tt.fields.exportProxyServiceInfoToExcelClient,
			}
			if got := w.ConnectivityCheckOfProxyService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConnectivityCheckOfProxyService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerStatisticsTransport_ExportProxyServiceInfoToExcel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]byte]](ctrl)

	type fields struct {
		getProxyServiceInfoClient             httpclientv1.ClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]]
		connectivityCheckOfProxyServiceClient httpclientv1.ClientBuilder[metav1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]]
		exportProxyServiceInfoToExcelClient   httpclientv1.ClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]byte]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]byte]]
	}{
		{
			name: "returns export proxy service info to excel client",
			fields: fields{
				exportProxyServiceInfoToExcelClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerStatisticsTransport{
				getProxyServiceInfoClient:             tt.fields.getProxyServiceInfoClient,
				connectivityCheckOfProxyServiceClient: tt.fields.connectivityCheckOfProxyServiceClient,
				exportProxyServiceInfoToExcelClient:   tt.fields.exportProxyServiceInfoToExcelClient,
			}
			if got := w.ExportProxyServiceInfoToExcel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExportProxyServiceInfoToExcel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerStatisticsTransport_GetProxyServiceInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]](ctrl)

	type fields struct {
		getProxyServiceInfoClient             httpclientv1.ClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]]
		connectivityCheckOfProxyServiceClient httpclientv1.ClientBuilder[metav1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]]
		exportProxyServiceInfoToExcelClient   httpclientv1.ClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]byte]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]]
	}{
		{
			name: "returns get proxy service info client",
			fields: fields{
				getProxyServiceInfoClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerStatisticsTransport{
				getProxyServiceInfoClient:             tt.fields.getProxyServiceInfoClient,
				connectivityCheckOfProxyServiceClient: tt.fields.connectivityCheckOfProxyServiceClient,
				exportProxyServiceInfoToExcelClient:   tt.fields.exportProxyServiceInfoToExcelClient,
			}
			if got := w.GetProxyServiceInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProxyServiceInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
