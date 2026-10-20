package service

import (
	"context"
	"encoding/json"
	"testing"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/marmotedu/errors"
	"github.com/tanganyu1114/heimdallr-reborn/server/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	epclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/endpoint"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"
	"go.uber.org/mock/gomock"
)

func Test_newApiService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockApiEndpoints(ctrl)
	mockEndpointFactory := epclientv1.NewMockFactory(ctrl)
	mockEndpointFactory.EXPECT().Apis().Return(mockEndpoints).AnyTimes()

	ctx := context.Background()
	svcFactory := &factory{ctx: ctx, eps: mockEndpointFactory}

	type args struct {
		factory *factory
	}
	tests := []struct {
		name string
		args args
		want ApiService
	}{
		{
			name: "successful creation",
			args: args{factory: svcFactory},
			want: &apiService{ctx: ctx, eps: mockEndpoints},
		},
		{
			name: "nil factory returns nil",
			args: args{factory: nil},
			want: nil,
		},
		{
			name: "nil eps in factory returns nil",
			args: args{factory: &factory{ctx: ctx, eps: nil}},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newApiService(tt.args.factory)
			if got == nil && tt.want != nil {
				t.Errorf("newApiService() = nil, want non-nil")
			}
			if got != nil && tt.want == nil {
				t.Errorf("newApiService() = %v, want nil", got)
			}
		})
	}
}

func Test_apiService_CreateApi(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockApiEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return modelclientv1.ResponseBody[httpclientv1.NilBody]{}, nil
	})

	mockErrorEndpoint := httpclientv1.NewEndpoint[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New("endpoint call error")
	})

	type fields struct {
		ctx context.Context
		eps epclientv1.ApiEndpoints
	}
	type args struct {
		api *model.SysApi
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setupMock func()
		wantErr   bool
	}{
		{
			name: "successful create API",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				api: &model.SysApi{
					Path:        "/api/test",
					Description: "Test API",
					ApiGroup:    "test",
					Method:      "POST",
				},
			},
			setupMock: func() {
				mockEndpoints.EXPECT().CreateApi().Return(mockEndpoint).Times(1)
			},
			wantErr: false,
		},
		{
			name: "endpoint call error",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				api: &model.SysApi{
					Path:        "/api/error",
					Description: "Error API",
					ApiGroup:    "test",
					Method:      "POST",
				},
			},
			setupMock: func() {
				mockEndpoints.EXPECT().CreateApi().Return(mockErrorEndpoint).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			s := &apiService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			if err := s.CreateApi(tt.args.api); (err != nil) != tt.wantErr {
				t.Errorf("CreateApi() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_apiService_UpdateApi(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockApiEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return modelclientv1.ResponseBody[httpclientv1.NilBody]{}, nil
	})

	mockErrorEndpoint := httpclientv1.NewEndpoint[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New("endpoint call error")
	})

	type fields struct {
		ctx context.Context
		eps epclientv1.ApiEndpoints
	}
	type args struct {
		api *model.SysApi
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setupMock func()
		wantErr   bool
	}{
		{
			name: "successful update API",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				api: &model.SysApi{
					Path:        "/api/test",
					Description: "Updated API",
					ApiGroup:    "test",
					Method:      "PUT",
				},
			},
			setupMock: func() {
				mockEndpoints.EXPECT().UpdateApi().Return(mockEndpoint).Times(1)
			},
			wantErr: false,
		},
		{
			name: "endpoint call error",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				api: &model.SysApi{
					Path:        "/api/error",
					Description: "Error API",
					ApiGroup:    "test",
					Method:      "PUT",
				},
			},
			setupMock: func() {
				mockEndpoints.EXPECT().UpdateApi().Return(mockErrorEndpoint).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			s := &apiService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			if err := s.UpdateApi(tt.args.api); (err != nil) != tt.wantErr {
				t.Errorf("UpdateApi() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_apiService_GetAllApis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockApiEndpoints(ctrl)
	ctx := context.Background()

	expectedApis := []model.SysApi{
		{
			Path:        "/api/test1",
			Description: "Test API 1",
			ApiGroup:    "test",
			Method:      "POST",
		},
		{
			Path:        "/api/test2",
			Description: "Test API 2",
			ApiGroup:    "test",
			Method:      "GET",
		},
	}

	mockEndpoint := httpclientv1.NewEndpoint[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]](func(ctx context.Context, req interface{}) (interface{}, error) {
		respData := &response.SysAPIListResponse{
			Apis: expectedApis,
		}
		data, _ := json.Marshal(respData)
		return modelclientv1.ResponseBody[*response.SysAPIListResponse]{Data: data}, nil
	})

	mockErrorEndpoint := httpclientv1.NewEndpoint[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New("endpoint call error")
	})

	type fields struct {
		ctx context.Context
		eps epclientv1.ApiEndpoints
	}
	tests := []struct {
		name      string
		fields    fields
		setupMock func()
		want      *response.SysAPIListResponse
		wantErr   bool
	}{
		{
			name: "successful get all APIs",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			setupMock: func() {
				mockEndpoints.EXPECT().GetAllApis().Return(mockEndpoint).Times(1)
			},
			want: &response.SysAPIListResponse{
				Apis: expectedApis,
			},
			wantErr: false,
		},
		{
			name: "endpoint call error",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			setupMock: func() {
				mockEndpoints.EXPECT().GetAllApis().Return(mockErrorEndpoint).Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			s := &apiService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.GetAllApis()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllApis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("GetAllApis() got = nil, want non-nil")
				return
			}
			if !tt.wantErr && len(got.Apis) != len(tt.want.Apis) {
				t.Errorf("GetAllApis() got %d APIs, want %d APIs", len(got.Apis), len(tt.want.Apis))
			}
		})
	}
}
