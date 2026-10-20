package auth

import (
	"testing"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/stretchr/testify/assert"
	"github.com/tanganyu1114/heimdallr-reborn/server/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/transport"
	"go.uber.org/mock/gomock"
)

func Test_apiMiddleware_CreateApi(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}{
		{
			name: "create API client builder with valid auth middleware",
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
			mockApiTransport := transport.NewMockApiTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)

			mockFactory.EXPECT().Apis().Return(mockApiTransport)
			mockApiTransport.EXPECT().CreateApi().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder).AnyTimes()
			mockClientBuilder.EXPECT().WithOptions(gomock.Any(), gomock.Any()).Return(mockClientBuilder).AnyTimes()

			tt.fields.md.txp = mockFactory

			a := &apiMiddleware{
				md: tt.fields.md,
			}
			got := a.CreateApi()
			assert.NotNil(t, got, "CreateApi() should return non-nil ClientBuilder")
		})
	}
}

func Test_apiMiddleware_GetAllApis(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]]
	}{
		{
			name: "get all APIs client builder with valid auth middleware",
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
			mockApiTransport := transport.NewMockApiTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]](ctrl)

			mockFactory.EXPECT().Apis().Return(mockApiTransport)
			mockApiTransport.EXPECT().GetAllApis().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder).AnyTimes()
			mockClientBuilder.EXPECT().WithOptions(gomock.Any(), gomock.Any()).Return(mockClientBuilder).AnyTimes()

			tt.fields.md.txp = mockFactory

			a := &apiMiddleware{
				md: tt.fields.md,
			}
			got := a.GetAllApis()
			assert.NotNil(t, got, "GetAllApis() should return non-nil ClientBuilder")
		})
	}
}

func Test_apiMiddleware_UpdateApi(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}{
		{
			name: "update API client builder with valid auth middleware",
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
			mockApiTransport := transport.NewMockApiTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)

			mockFactory.EXPECT().Apis().Return(mockApiTransport)
			mockApiTransport.EXPECT().UpdateApi().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder).AnyTimes()
			mockClientBuilder.EXPECT().WithOptions(gomock.Any(), gomock.Any()).Return(mockClientBuilder).AnyTimes()

			tt.fields.md.txp = mockFactory

			a := &apiMiddleware{
				md: tt.fields.md,
			}
			got := a.UpdateApi()
			assert.NotNil(t, got, "UpdateApi() should return non-nil ClientBuilder")
		})
	}
}

func Test_newApisMiddleware(t *testing.T) {
	type args struct {
		md *authMiddleware
	}
	tests := []struct {
		name string
		args args
		want *apiMiddleware
	}{
		{
			name: "create api middleware with valid auth middleware",
			args: args{
				md: &authMiddleware{
					apiKey:     "test-api-key",
					apiSecret:  "test-api-secret",
					bufferTime: 3600,
					token:      "test-token",
					expiresAt:  1000000,
				},
			},
			want: &apiMiddleware{
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
			name: "create api middleware with empty credentials",
			args: args{
				md: &authMiddleware{
					apiKey:     "",
					apiSecret:  "",
					bufferTime: 0,
					token:      "",
					expiresAt:  0,
				},
			},
			want: &apiMiddleware{
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
			got := newApisMiddleware(tt.args.md)
			assert.NotNil(t, got, "newApisMiddleware() should return non-nil apiMiddleware")
			assert.Equal(t, tt.want.md, got.md, "newApisMiddleware() should set md field correctly")
		})
	}
}
