package endpoint

import (
	"testing"

	"github.com/tanganyu1114/heimdallr-reborn/server/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/transport"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_newApiEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockFactory(ctrl)
	mockApiTransport := transport.NewMockApiTransport(ctrl)
	mockTransport.EXPECT().Apis().Return(mockApiTransport).AnyTimes()

	type args struct {
		f *factory
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "creates api endpoints",
			args: args{
				f: &factory{
					transport: mockTransport,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newApiEndpoints(tt.args.f)
			if got == nil {
				t.Errorf("newApiEndpoints() = nil, want non-nil")
			}
		})
	}
}

func Test_apiEndpoints_CreateApi(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockApiTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockClient := httpclientv1.NewMockClient[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().CreateApi().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.ApiTransport
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "returns create API endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiEndpoints{
				transport: tt.fields.transport,
			}
			got := a.CreateApi()
			if got == nil {
				t.Errorf("CreateApi() = nil, want non-nil")
			}
		})
	}
}

func Test_apiEndpoints_UpdateApi(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockApiTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockClient := httpclientv1.NewMockClient[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().UpdateApi().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.ApiTransport
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "returns update API endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiEndpoints{
				transport: tt.fields.transport,
			}
			got := a.UpdateApi()
			if got == nil {
				t.Errorf("UpdateApi() = nil, want non-nil")
			}
		})
	}
}

func Test_apiEndpoints_GetAllApis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockApiTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]](ctrl)
	mockClient := httpclientv1.NewMockClient[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().GetAllApis().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.ApiTransport
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "returns get all APIs endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiEndpoints{
				transport: tt.fields.transport,
			}
			got := a.GetAllApis()
			if got == nil {
				t.Errorf("GetAllApis() = nil, want non-nil")
			}
		})
	}
}
