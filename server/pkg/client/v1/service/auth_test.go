package service

import (
	"context"
	"encoding/json"
	"testing"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/marmotedu/errors"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	epclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/endpoint"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"
	"go.uber.org/mock/gomock"
)

func Test_newAuthService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockAuthEndpoints(ctrl)
	mockEndpointFactory := epclientv1.NewMockFactory(ctrl)
	mockEndpointFactory.EXPECT().Authz().Return(mockEndpoints).AnyTimes()

	ctx := context.Background()
	svcFactory := &factory{ctx: ctx, eps: mockEndpointFactory}

	type args struct {
		factory *factory
	}
	tests := []struct {
		name string
		args args
		want AuthService
	}{
		{
			name: "successful creation",
			args: args{factory: svcFactory},
			want: &authService{ctx: ctx, eps: mockEndpoints},
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
			got := newAuthService(tt.args.factory)
			if got == nil && tt.want != nil {
				t.Errorf("newAuthService() = nil, want non-nil")
			}
			if got != nil && tt.want == nil {
				t.Errorf("newAuthService() = %v, want nil", got)
			}
		})
	}
}

func Test_authService_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockAuthEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]](func(ctx context.Context, req interface{}) (interface{}, error) {
		claims := &response.ValidateTokenResponse{
			Claims: &request.CustomClaims{
				Username: "testuser",
			},
		}
		data, _ := json.Marshal(claims)
		return modelclientv1.ResponseBody[*response.ValidateTokenResponse]{Data: data}, nil
	})

	mockErrorEndpoint := httpclientv1.NewEndpoint[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New("endpoint call error")
	})

	type fields struct {
		ctx context.Context
		eps epclientv1.AuthEndpoints
	}
	type args struct {
		token string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setupMock func()
		want      *response.ValidateTokenResponse
		wantErr   bool
	}{
		{
			name: "successful validate token",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				token: "valid-token",
			},
			setupMock: func() {
				mockEndpoints.EXPECT().ValidateToken().Return(mockEndpoint).Times(1)
			},
			want: &response.ValidateTokenResponse{
				Claims: &request.CustomClaims{
					Username: "testuser",
				},
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
				token: "invalid-token",
			},
			setupMock: func() {
				mockEndpoints.EXPECT().ValidateToken().Return(mockErrorEndpoint).Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			s := &authService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.ValidateToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want != nil {
				t.Errorf("ValidateToken() got = nil, want %v", tt.want)
			}
		})
	}
}

func Test_authService_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockAuthEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]](func(ctx context.Context, req interface{}) (interface{}, error) {
		respData := &response.RefreshTokenResponse{
			Token:     "new-refreshed-token",
			ExpiresAt: 1234567890,
		}
		data, _ := json.Marshal(respData)
		return modelclientv1.ResponseBody[*response.RefreshTokenResponse]{Data: data}, nil
	})

	mockErrorEndpoint := httpclientv1.NewEndpoint[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New("endpoint call error")
	})

	type fields struct {
		ctx context.Context
		eps epclientv1.AuthEndpoints
	}
	type args struct {
		token     string
		expiresAt int64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setupMock func()
		want      *response.RefreshTokenResponse
		wantErr   bool
	}{
		{
			name: "successful refresh token",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				token:     "old-token",
				expiresAt: 1234567890,
			},
			setupMock: func() {
				mockEndpoints.EXPECT().RefreshToken().Return(mockEndpoint).Times(1)
			},
			want: &response.RefreshTokenResponse{
				Token:     "new-refreshed-token",
				ExpiresAt: 1234567890,
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
				token:     "expired-token",
				expiresAt: 0,
			},
			setupMock: func() {
				mockEndpoints.EXPECT().RefreshToken().Return(mockErrorEndpoint).Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			s := &authService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.RefreshToken(tt.args.token, tt.args.expiresAt)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want != nil {
				t.Errorf("RefreshToken() got = nil, want %v", tt.want)
			}
		})
	}
}

func Test_authService_CheckPermission(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockAuthEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]](func(ctx context.Context, req interface{}) (interface{}, error) {
		respData := &response.CheckPermissionResponse{
			Allowed: true,
		}
		data, _ := json.Marshal(respData)
		return modelclientv1.ResponseBody[*response.CheckPermissionResponse]{Data: data}, nil
	})

	mockErrorEndpoint := httpclientv1.NewEndpoint[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New("endpoint call error")
	})

	type fields struct {
		ctx context.Context
		eps epclientv1.AuthEndpoints
	}
	type args struct {
		sub string
		obj string
		act string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setupMock func()
		want      *response.CheckPermissionResponse
		wantErr   bool
	}{
		{
			name: "successful check permission allowed",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				sub: "user-123",
				obj: "/api/data",
				act: "GET",
			},
			setupMock: func() {
				mockEndpoints.EXPECT().CheckPermission().Return(mockEndpoint).Times(1)
			},
			want: &response.CheckPermissionResponse{
				Allowed: true,
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
				sub: "user-456",
				obj: "/api/admin",
				act: "DELETE",
			},
			setupMock: func() {
				mockEndpoints.EXPECT().CheckPermission().Return(mockErrorEndpoint).Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			s := &authService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.CheckPermission(tt.args.sub, tt.args.obj, tt.args.act)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPermission() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want != nil {
				t.Errorf("CheckPermission() got = nil, want %v", tt.want)
			}
		})
	}
}

func Test_authService_LogOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockAuthEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return modelclientv1.ResponseBody[httpclientv1.NilBody]{}, nil
	})

	mockErrorEndpoint := httpclientv1.NewEndpoint[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New("endpoint call error")
	})

	type fields struct {
		ctx context.Context
		eps epclientv1.AuthEndpoints
	}
	type args struct {
		record *request.LogOperationRequest
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setupMock func()
		wantErr   bool
	}{
		{
			name: "successful log operation",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				record: &request.LogOperationRequest{
					IP:         "127.0.0.1",
					Method:     "POST",
					Path:       "/api/data",
					Status:     200,
					Latency:    1000000,
					Agent:      "Mozilla/5.0",
					Body:       `{"key":"value"}`,
					Resp:       `{"code":0}`,
					UserID:     1,
					IsSDKLogin: false,
				},
			},
			setupMock: func() {
				mockEndpoints.EXPECT().LogOperation().Return(mockEndpoint).Times(1)
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
				record: &request.LogOperationRequest{
					IP:         "127.0.0.1",
					Method:     "POST",
					Path:       "/api/error",
					Status:     500,
					Latency:    1000000,
					Agent:      "Mozilla/5.0",
					Body:       `{}`,
					Resp:       `{}`,
					UserID:     1,
					IsSDKLogin: false,
				},
			},
			setupMock: func() {
				mockEndpoints.EXPECT().LogOperation().Return(mockErrorEndpoint).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			s := &authService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			if err := s.LogOperation(tt.args.record); (err != nil) != tt.wantErr {
				t.Errorf("LogOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
