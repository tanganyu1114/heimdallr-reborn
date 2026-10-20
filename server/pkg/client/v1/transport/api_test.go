package transport

import (
	"reflect"
	"testing"

	"github.com/tanganyu1114/heimdallr-reborn/server/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_newApiTransport(t *testing.T) {
	type args struct {
		transport *transport
	}
	tests := []struct {
		name string
		args args
		want ApiTransport
	}{
		{
			name: "creates api transport with valid base URL",
			args: args{
				transport: &transport{
					baseURL: "http://localhost:8080",
				},
			},
		},
		{
			name: "creates api transport with empty base URL",
			args: args{
				transport: &transport{
					baseURL: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newApiTransport(tt.args.transport)
			if got == nil {
				t.Errorf("newApiTransport() = nil, want non-nil")
			}
		})
	}
}

func Test_apiTransport_CreateApi(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)

	type fields struct {
		createApiClient  httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
		updateApiClient  httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
		getAllApisClient httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}{
		{
			name: "returns create API client",
			fields: fields{
				createApiClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiTransport{
				createApiClient:  tt.fields.createApiClient,
				updateApiClient:  tt.fields.updateApiClient,
				getAllApisClient: tt.fields.getAllApisClient,
			}
			if got := a.CreateApi(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateApi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiTransport_UpdateApi(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)

	type fields struct {
		createApiClient  httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
		updateApiClient  httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
		getAllApisClient httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}{
		{
			name: "returns update API client",
			fields: fields{
				updateApiClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiTransport{
				createApiClient:  tt.fields.createApiClient,
				updateApiClient:  tt.fields.updateApiClient,
				getAllApisClient: tt.fields.getAllApisClient,
			}
			if got := a.UpdateApi(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateApi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiTransport_GetAllApis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]](ctrl)

	type fields struct {
		createApiClient  httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
		updateApiClient  httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
		getAllApisClient httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]]
	}{
		{
			name: "returns get all APIs client",
			fields: fields{
				getAllApisClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiTransport{
				createApiClient:  tt.fields.createApiClient,
				updateApiClient:  tt.fields.updateApiClient,
				getAllApisClient: tt.fields.getAllApisClient,
			}
			if got := a.GetAllApis(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllApis() = %v, want %v", got, tt.want)
			}
		})
	}
}
