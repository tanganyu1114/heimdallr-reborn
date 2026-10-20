package auth

import (
	"testing"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/stretchr/testify/assert"
	"github.com/tanganyu1114/heimdallr-reborn/server/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/transport"
	"go.uber.org/mock/gomock"
)

func Test_authsMiddleware_CheckPermission(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}{
		{
			name: "check permission client builder with valid auth middleware",
			fields: fields{
				md: &authMiddleware{
					apiKey:     "test-api-key",
					apiSecret:  "test-api-secret",
					bufferTime: 3600,
					token:      "test-token",
					expiresAt:  1000000,
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockAuthzTransport := transport.NewMockAuthTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]](ctrl)

			mockFactory.EXPECT().Authz().Return(mockAuthzTransport)
			mockAuthzTransport.EXPECT().CheckPermission().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder).AnyTimes()
			mockClientBuilder.EXPECT().WithOptions(gomock.Any(), gomock.Any()).Return(mockClientBuilder).AnyTimes()

			tt.fields.md.txp = mockFactory

			a := &authsMiddleware{
				md: tt.fields.md,
			}
			got := a.CheckPermission()
			assert.NotNil(t, got, "CheckPermission() should return non-nil ClientBuilder")
		})
	}
}

func Test_authsMiddleware_LogOperation(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}{
		{
			name: "log operation client builder with valid auth middleware",
			fields: fields{
				md: &authMiddleware{
					apiKey:     "test-api-key",
					apiSecret:  "test-api-secret",
					bufferTime: 3600,
					token:      "test-token",
					expiresAt:  1000000,
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockAuthzTransport := transport.NewMockAuthTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)

			mockFactory.EXPECT().Authz().Return(mockAuthzTransport)
			mockAuthzTransport.EXPECT().LogOperation().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder).AnyTimes()
			mockClientBuilder.EXPECT().WithOptions(gomock.Any(), gomock.Any()).Return(mockClientBuilder).AnyTimes()

			tt.fields.md.txp = mockFactory

			a := &authsMiddleware{
				md: tt.fields.md,
			}
			got := a.LogOperation()
			assert.NotNil(t, got, "LogOperation() should return non-nil ClientBuilder")
		})
	}
}

func Test_authsMiddleware_RefreshToken(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]]
	}{
		{
			name: "refresh token client builder with valid auth middleware",
			fields: fields{
				md: &authMiddleware{
					apiKey:     "test-api-key",
					apiSecret:  "test-api-secret",
					bufferTime: 3600,
					token:      "test-token",
					expiresAt:  1000000,
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockAuthzTransport := transport.NewMockAuthTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]](ctrl)

			mockFactory.EXPECT().Authz().Return(mockAuthzTransport)
			mockAuthzTransport.EXPECT().RefreshToken().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder).AnyTimes()
			mockClientBuilder.EXPECT().WithOptions(gomock.Any(), gomock.Any()).Return(mockClientBuilder).AnyTimes()

			tt.fields.md.txp = mockFactory

			a := &authsMiddleware{
				md: tt.fields.md,
			}
			got := a.RefreshToken()
			assert.NotNil(t, got, "RefreshToken() should return non-nil ClientBuilder")
		})
	}
}

func Test_authsMiddleware_ValidateToken(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]]
	}{
		{
			name: "validate token client builder with valid auth middleware",
			fields: fields{
				md: &authMiddleware{
					apiKey:     "test-api-key",
					apiSecret:  "test-api-secret",
					bufferTime: 3600,
					token:      "test-token",
					expiresAt:  1000000,
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockAuthzTransport := transport.NewMockAuthTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]](ctrl)

			mockFactory.EXPECT().Authz().Return(mockAuthzTransport)
			mockAuthzTransport.EXPECT().ValidateToken().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder).AnyTimes()
			mockClientBuilder.EXPECT().WithOptions(gomock.Any(), gomock.Any()).Return(mockClientBuilder).AnyTimes()

			tt.fields.md.txp = mockFactory

			a := &authsMiddleware{
				md: tt.fields.md,
			}
			got := a.ValidateToken()
			assert.NotNil(t, got, "ValidateToken() should return non-nil ClientBuilder")
		})
	}
}

func Test_newAuthzMiddleware(t *testing.T) {
	type args struct {
		md *authMiddleware
	}
	tests := []struct {
		name string
		args args
		want *authsMiddleware
	}{
		{
			name: "create authz middleware with valid auth middleware",
			args: args{
				md: &authMiddleware{
					apiKey:     "test-api-key",
					apiSecret:  "test-api-secret",
					bufferTime: 3600,
					token:      "test-token",
					expiresAt:  1000000,
				},
			},
			want: &authsMiddleware{
				md: &authMiddleware{
					apiKey:     "test-api-key",
					apiSecret:  "test-api-secret",
					bufferTime: 3600,
					token:      "test-token",
					expiresAt:  1000000,
				},
			},
		},
		{
			name: "create authz middleware with empty credentials",
			args: args{
				md: &authMiddleware{
					apiKey:     "",
					apiSecret:  "",
					bufferTime: 0,
					token:      "",
					expiresAt:  0,
				},
			},
			want: &authsMiddleware{
				md: &authMiddleware{
					apiKey:     "",
					apiSecret:  "",
					bufferTime: 0,
					token:      "",
					expiresAt:  0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newAuthzMiddleware(tt.args.md)
			assert.NotNil(t, got, "newAuthzMiddleware() should return non-nil authsMiddleware")
			assert.Equal(t, tt.want.md, got.md, "newAuthzMiddleware() should set md field correctly")
		})
	}
}
