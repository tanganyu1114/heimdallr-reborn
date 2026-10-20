package endpoint

import (
	"testing"

	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/transport"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_newAuthEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockFactory(ctrl)
	mockAuthTransport := transport.NewMockAuthTransport(ctrl)
	mockTransport.EXPECT().Authz().Return(mockAuthTransport).AnyTimes()

	type args struct {
		f *factory
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "creates auth endpoints",
			args: args{
				f: &factory{
					transport: mockTransport,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newAuthEndpoints(tt.args.f)
			if got == nil {
				t.Errorf("newAuthEndpoints() = nil, want non-nil")
			}
		})
	}
}

func Test_authEndpoints_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockAuthTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]](ctrl)
	mockClient := httpclientv1.NewMockClient[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().ValidateToken().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.AuthTransport
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "returns validate token endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authEndpoints{
				transport: tt.fields.transport,
			}
			got := a.ValidateToken()
			if got == nil {
				t.Errorf("ValidateToken() = nil, want non-nil")
			}
		})
	}
}

func Test_authEndpoints_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockAuthTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]](ctrl)
	mockClient := httpclientv1.NewMockClient[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().RefreshToken().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.AuthTransport
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "returns refresh token endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authEndpoints{
				transport: tt.fields.transport,
			}
			got := a.RefreshToken()
			if got == nil {
				t.Errorf("RefreshToken() = nil, want non-nil")
			}
		})
	}
}

func Test_authEndpoints_CheckPermission(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockAuthTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]](ctrl)
	mockClient := httpclientv1.NewMockClient[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().CheckPermission().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.AuthTransport
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "returns check permission endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authEndpoints{
				transport: tt.fields.transport,
			}
			got := a.CheckPermission()
			if got == nil {
				t.Errorf("CheckPermission() = nil, want non-nil")
			}
		})
	}
}

func Test_authEndpoints_LogOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockAuthTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockClient := httpclientv1.NewMockClient[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().LogOperation().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.AuthTransport
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "returns log operation endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authEndpoints{
				transport: tt.fields.transport,
			}
			got := a.LogOperation()
			if got == nil {
				t.Errorf("LogOperation() = nil, want non-nil")
			}
		})
	}
}
