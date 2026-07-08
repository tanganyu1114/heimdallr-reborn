package service

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	epclientv1 "gin-vue-admin/pkg/client/v1/endpoint"
	modelclientv1 "gin-vue-admin/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_newWebServerStatisticsService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerStatisticsEndpoints(ctrl)
	mockEndpointFactory := epclientv1.NewMockFactory(ctrl)
	mockEndpointFactory.EXPECT().WebServerStatistics().Return(mockEndpoints).AnyTimes()

	ctx := context.Background()
	svcFactory := &factory{ctx: ctx, eps: mockEndpointFactory}

	type args struct {
		factory *factory
	}
	tests := []struct {
		name string
		args args
		want WebServerStatisticsService
	}{
		{
			name: "successful creation",
			args: args{factory: svcFactory},
			want: &webServerStatisticsService{ctx: ctx, eps: mockEndpoints},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newWebServerStatisticsService(tt.args.factory)
			if !reflect.DeepEqual(reflect.TypeOf(got), reflect.TypeOf(tt.want)) {
				t.Errorf("newWebServerStatisticsService() type = %T, want %T", got, tt.want)
			}
		})
	}
}

func Test_webServerStatisticsService_ConnectivityCheckOfProxyService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerStatisticsEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[metav1.ConnectivityCheckOfProxiedServersRequestOptions, modelclientv1.ResponseBody[v1.ProxyServiceInfo]](func(ctx context.Context, req interface{}) (interface{}, error) {
		data, _ := json.Marshal(v1.ProxyServiceInfo{})
		return modelclientv1.ResponseBody[v1.ProxyServiceInfo]{Data: data}, nil
	})
	mockEndpoints.EXPECT().ConnectivityCheckOfProxyService().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerStatisticsEndpoints
	}
	type args struct {
		opts *metav1.ConnectivityCheckOfProxiedServersRequestOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "successful connectivity check",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &metav1.ConnectivityCheckOfProxiedServersRequestOptions{},
			},
			want:    v1.ProxyServiceInfo{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerStatisticsService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.ConnectivityCheckOfProxyService(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConnectivityCheckOfProxyService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConnectivityCheckOfProxyService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerStatisticsService_ExportProxyServiceInfoToExcel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerStatisticsEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[[]byte]](func(ctx context.Context, req interface{}) (interface{}, error) {
		data, _ := json.Marshal([]byte("excel-data"))
		return modelclientv1.ResponseBody[[]byte]{Data: data}, nil
	})
	mockEndpoints.EXPECT().ExportProxyServiceInfoToExcel().Return(mockEndpoint).AnyTimes()

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerStatisticsEndpoints
	}
	type args struct {
		opts *metav1.WebServerOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "successful export proxy service info to excel",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &metav1.WebServerOptions{},
			},
			want:    []byte("excel-data"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerStatisticsService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.ExportProxyServiceInfoToExcel(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExportProxyServiceInfoToExcel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("ExportProxyServiceInfoToExcel() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerStatisticsService_GetProxyServiceInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerStatisticsEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerOptions, modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]](func(ctx context.Context, req interface{}) (interface{}, error) {
		data, _ := json.Marshal([]v1.ProxyServiceInfo{{ServerName: "test-proxy"}})
		return modelclientv1.ResponseBody[[]v1.ProxyServiceInfo]{Data: data}, nil
	})
	mockEndpoints.EXPECT().GetProxyServiceInfo().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerStatisticsEndpoints
	}
	type args struct {
		opts *metav1.WebServerOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []v1.ProxyServiceInfo
		wantErr bool
	}{
		{
			name: "successful get proxy service info",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &metav1.WebServerOptions{},
			},
			want:    []v1.ProxyServiceInfo{{ServerName: "test-proxy"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerStatisticsService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.GetProxyServiceInfo(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProxyServiceInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetProxyServiceInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
