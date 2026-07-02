package auth

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/model/request"
	"gin-vue-admin/model/response"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"gin-vue-admin/pkg/client/v1/transport"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/marmotedu/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthMiddleware_GetToken(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
		token     string
		expiresAt int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "empty token",
			fields: fields{apiKey: "key", apiSecret: "secret", token: "", expiresAt: 0},
			want:   "",
		},
		{
			name:   "valid token",
			fields: fields{apiKey: "key", apiSecret: "secret", token: "test-token", expiresAt: time.Now().Unix() + 3600},
			want:   "test-token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, 3600)
			authMw := mw(mockFactory).(*authMiddleware)

			authMw.mu.Lock()
			authMw.token = tt.fields.token
			authMw.expiresAt = tt.fields.expiresAt
			authMw.mu.Unlock()

			assert.Equal(t, tt.want, authMw.GetToken())
		})
	}
}

func TestAuthMiddleware_GetExpiresAt(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
		expiresAt int64
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name:   "zero expires at",
			fields: fields{apiKey: "key", apiSecret: "secret", expiresAt: 0},
			want:   0,
		},
		{
			name:   "future expires at",
			fields: fields{apiKey: "key", apiSecret: "secret", expiresAt: 1000000},
			want:   1000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, 3600)
			authMw := mw(mockFactory).(*authMiddleware)

			authMw.mu.Lock()
			authMw.expiresAt = tt.fields.expiresAt
			authMw.mu.Unlock()

			assert.Equal(t, tt.want, authMw.GetExpiresAt())
		})
	}
}

func TestAuthMiddleware_GetAPIKey(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "get api key",
			fields: fields{apiKey: "test-api-key", apiSecret: "test-api-secret"},
			want:   "test-api-key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, 3600)
			authMw := mw(mockFactory).(*authMiddleware)

			assert.Equal(t, tt.want, authMw.GetAPIKey())
		})
	}
}

func TestAuthMiddleware_GetAPISecret(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "get api secret",
			fields: fields{apiKey: "test-api-key", apiSecret: "test-api-secret"},
			want:   "test-api-secret",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, 3600)
			authMw := mw(mockFactory).(*authMiddleware)

			assert.Equal(t, tt.want, authMw.GetAPISecret())
		})
	}
}

func TestAuthMiddleware_SysUsers(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "get sys users transport",
			fields: fields{apiKey: "key", apiSecret: "secret"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
			mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)

			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, 3600)
			authMw := mw(mockFactory).(*authMiddleware)

			got := authMw.SysUsers()
			assert.Equal(t, mockSysUserTransport, got)
		})
	}
}

func TestAuthMiddleware_AgentInfos(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "get agent infos transport",
			fields: fields{apiKey: "key", apiSecret: "secret"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, 3600)
			authMw := mw(mockFactory).(*authMiddleware)

			got := authMw.AgentInfos()
			assert.NotNil(t, got)
		})
	}
}

func TestAuthMiddleware_Groups(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "get groups transport",
			fields: fields{apiKey: "key", apiSecret: "secret"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, 3600)
			authMw := mw(mockFactory).(*authMiddleware)

			got := authMw.Groups()
			assert.NotNil(t, got)
		})
	}
}

func TestAuthMiddleware_Hosts(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "get hosts transport",
			fields: fields{apiKey: "key", apiSecret: "secret"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, 3600)
			authMw := mw(mockFactory).(*authMiddleware)

			got := authMw.Hosts()
			assert.NotNil(t, got)
		})
	}
}

func TestAuthMiddleware_WebServerConfigs(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "get web server configs transport",
			fields: fields{apiKey: "key", apiSecret: "secret"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, 3600)
			authMw := mw(mockFactory).(*authMiddleware)

			got := authMw.WebServerConfigs()
			assert.NotNil(t, got)
		})
	}
}

func TestAuthMiddleware_WebServerBinCMDs(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "get web server bin cmds transport",
			fields: fields{apiKey: "key", apiSecret: "secret"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, 3600)
			authMw := mw(mockFactory).(*authMiddleware)

			got := authMw.WebServerBinCMDs()
			assert.NotNil(t, got)
		})
	}
}

func TestAuthMiddleware_WebServerStatistics(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "get web server statistics transport",
			fields: fields{apiKey: "key", apiSecret: "secret"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, 3600)
			authMw := mw(mockFactory).(*authMiddleware)

			got := authMw.WebServerStatistics()
			assert.NotNil(t, got)
		})
	}
}

func TestAuthMiddleware_EnsureValidToken(t *testing.T) {
	type fields struct {
		apiKey     string
		apiSecret  string
		bufferTime int64
		token      string
		expiresAt  int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errMsg  string
	}{
		{
			name: "no refresh needed - valid token",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
				token:     "valid-token",
				expiresAt: time.Now().Unix() + 7200,
			},
			wantErr: false,
		},
		{
			name: "token not exists",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
				token:     "",
				expiresAt: 0,
			},
			wantErr: true,
			errMsg:  "sdkLogin client builder is nil",
		},
		{
			name: "token expiring soon",
			fields: fields{
				apiKey:     "key",
				apiSecret:  "secret",
				bufferTime: 3600,
				token:      "old-token",
				expiresAt:  time.Now().Unix() + 100,
			},
			wantErr: true,
			errMsg:  "sdkLogin client builder is nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)

			if tt.fields.token == "" || (tt.fields.expiresAt > 0 && tt.fields.expiresAt-time.Now().Unix() < tt.fields.bufferTime) {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport).AnyTimes()
				mockSysUserTransport.EXPECT().SDKLogin().Return(nil).AnyTimes()
			}

			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
			authMw := mw(mockFactory).(*authMiddleware)

			authMw.mu.Lock()
			authMw.token = tt.fields.token
			authMw.expiresAt = tt.fields.expiresAt
			authMw.mu.Unlock()

			err := authMw.ensureValidToken()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.fields.token, authMw.GetToken())
			}
		})
	}
}

func TestAuthMiddleware_PassiveRefreshToken(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
		token     string
		expiresAt int64
	}
	type args struct {
		ctx      context.Context
		response *http.Response
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantToken    string
		wantExpires  int64
		wantCtxToken string
	}{
		{
			name:   "success - new token in response",
			fields: fields{apiKey: "key", apiSecret: "secret", token: "old-token", expiresAt: 1000000},
			args: args{
				ctx: context.Background(),
				response: &http.Response{
					Header: http.Header{
						"New-Token":      []string{"new-token"},
						"New-Expires-At": []string{strconv.FormatInt(2000000, 10)},
					},
				},
			},
			wantToken:    "new-token",
			wantExpires:  2000000,
			wantCtxToken: "new-token",
		},
		{
			name:   "no new token in response",
			fields: fields{apiKey: "key", apiSecret: "secret", token: "old-token", expiresAt: 2000000},
			args: args{
				ctx:      context.Background(),
				response: &http.Response{Header: http.Header{}},
			},
			wantToken:    "old-token",
			wantExpires:  2000000,
			wantCtxToken: "",
		},
		{
			name:   "invalid expires at",
			fields: fields{apiKey: "key", apiSecret: "secret", token: "old-token", expiresAt: 2000000},
			args: args{
				ctx: context.Background(),
				response: &http.Response{
					Header: http.Header{
						"new-token":      []string{"new-token"},
						"new-expires-at": []string{"invalid"},
					},
				},
			},
			wantToken:    "old-token",
			wantExpires:  2000000,
			wantCtxToken: "",
		},
		{
			name:   "older token in response",
			fields: fields{apiKey: "key", apiSecret: "secret", token: "old-token", expiresAt: 2000000},
			args: args{
				ctx: context.Background(),
				response: &http.Response{
					Header: http.Header{
						"new-token":      []string{"new-token"},
						"new-expires-at": []string{strconv.FormatInt(1000000, 10)},
					},
				},
			},
			wantToken:    "old-token",
			wantExpires:  2000000,
			wantCtxToken: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, 3600)
			authMw := mw(mockFactory).(*authMiddleware)

			authMw.mu.Lock()
			authMw.token = tt.fields.token
			authMw.expiresAt = tt.fields.expiresAt
			authMw.mu.Unlock()

			newCtx := authMw.passiveRefreshToken(tt.args.ctx, tt.args.response)

			assert.Equal(t, tt.wantToken, authMw.GetToken())
			assert.Equal(t, tt.wantExpires, authMw.GetExpiresAt())
			if tt.wantCtxToken != "" {
				gotToken, ok := newCtx.Value(RequestTokenKey).(string)
				assert.True(t, ok)
				assert.Equal(t, tt.wantCtxToken, gotToken)
			} else {
				assert.Equal(t, tt.args.ctx, newCtx)
			}
		})
	}
}

func TestAuthMiddleware_ConcurrentAccess(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
		token     string
		expiresAt int64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "concurrent read access",
			fields: fields{apiKey: "key", apiSecret: "secret", token: "test-token", expiresAt: time.Now().Unix() + 7200},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, 3600)
			authMw := mw(mockFactory).(*authMiddleware)

			authMw.mu.Lock()
			authMw.token = tt.fields.token
			authMw.expiresAt = tt.fields.expiresAt
			authMw.mu.Unlock()

			var wg sync.WaitGroup
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_ = authMw.GetToken()
					_ = authMw.GetExpiresAt()
					_ = authMw.GetAPIKey()
					_ = authMw.GetAPISecret()
				}()
			}
			wg.Wait()
		})
	}
}

func Test_authMiddleware_authOptions(t *testing.T) {
	type fields struct {
		apiKey     string
		apiSecret  string
		bufferTime int64
		token      string
		expiresAt  int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantLen int
	}{
		{
			name:    "returns client options",
			fields:  fields{apiKey: "key", apiSecret: "secret", bufferTime: 3600, token: "test-token", expiresAt: time.Now().Unix() + 7200},
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
			authMw := mw(mockFactory).(*authMiddleware)

			authMw.mu.Lock()
			authMw.token = tt.fields.token
			authMw.expiresAt = tt.fields.expiresAt
			authMw.mu.Unlock()

			opts := authMw.authOptions()
			assert.Len(t, opts, tt.wantLen)
		})
	}
}

func Test_authMiddleware_ensureValidTokenToCtx(t *testing.T) {
	type fields struct {
		apiKey     string
		apiSecret  string
		bufferTime int64
		token      string
		expiresAt  int64
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantToken string
		wantErr   bool
	}{
		{
			name: "valid token - add to context",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
				token:     "test-token",
				expiresAt: time.Now().Unix() + 7200,
			},
			args:      args{ctx: context.Background()},
			wantToken: "test-token",
			wantErr:   false,
		},
		{
			name: "token not exists - refresh fails",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
				token:     "",
				expiresAt: 0,
			},
			args:      args{ctx: context.Background()},
			wantToken: "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)

			if tt.fields.token == "" || (tt.fields.expiresAt > 0 && tt.fields.expiresAt-time.Now().Unix() < tt.fields.bufferTime) {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport).AnyTimes()
				mockSysUserTransport.EXPECT().SDKLogin().Return(nil).AnyTimes()
			}

			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
			authMw := mw(mockFactory).(*authMiddleware)

			authMw.mu.Lock()
			authMw.token = tt.fields.token
			authMw.expiresAt = tt.fields.expiresAt
			authMw.mu.Unlock()

			got, err := authMw.ensureValidTokenToCtx(tt.args.ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				token, ok := got.Value(RequestTokenKey).(string)
				assert.True(t, ok)
				assert.Equal(t, tt.wantToken, token)
			}
		})
	}
}

func Test_authMiddleware_refreshToken(t *testing.T) {
	type fields struct {
		apiKey     string
		apiSecret  string
		bufferTime int64
		token      string
		expiresAt  int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errMsg  string
	}{
		{
			name: "transport is nil",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
			},
			wantErr: true,
			errMsg:  "transport is nil",
		},
		{
			name: "sysUser transport is nil",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
			},
			wantErr: true,
			errMsg:  "sysUser transport is nil",
		},
		{
			name: "sdkLogin client builder is nil",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
			},
			wantErr: true,
			errMsg:  "sdkLogin client builder is nil",
		},
		{
			name: "sdkLogin client is nil",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
			},
			wantErr: true,
			errMsg:  "sdkLogin client is nil",
		},
		{
			name: "sdk login endpoint is nil",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
			},
			wantErr: true,
			errMsg:  "sdk login endpoint is nil",
		},
		{
			name: "login request fails",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
			},
			wantErr: true,
			errMsg:  "login error",
		},
		{
			name: "login success",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)

			if tt.errMsg == "transport is nil" {
				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)
				authMw.txp = nil

				err := authMw.refreshToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "sysUser transport is nil" {
				mockFactory.EXPECT().SysUsers().Return(nil)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.refreshToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "sdkLogin client builder is nil" {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().SDKLogin().Return(nil)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.refreshToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "sdkLogin client is nil" {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockClientBuilder := httpclientv1.NewMockClientBuilder[*request.SDKLogin, *response.LoginResponse](ctrl)
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().SDKLogin().Return(mockClientBuilder)
				mockClientBuilder.EXPECT().Build().Return(nil)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.refreshToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "sdk login endpoint is nil" {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockClientBuilder := httpclientv1.NewMockClientBuilder[*request.SDKLogin, *response.LoginResponse](ctrl)
				mockClient := httpclientv1.NewMockClient[*request.SDKLogin, *response.LoginResponse](ctrl)
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().SDKLogin().Return(mockClientBuilder)
				mockClientBuilder.EXPECT().Build().Return(mockClient)
				mockClient.EXPECT().Endpoint().Return(nil)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.refreshToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "login error" {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockClientBuilder := httpclientv1.NewMockClientBuilder[*request.SDKLogin, *response.LoginResponse](ctrl)
				mockClient := httpclientv1.NewMockClient[*request.SDKLogin, *response.LoginResponse](ctrl)
				mockEndpoint := func(ctx context.Context, req httpclientv1.HTTPRequest[*request.SDKLogin]) (*response.LoginResponse, error) {
					return nil, errors.New("login error")
				}
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().SDKLogin().Return(mockClientBuilder)
				mockClientBuilder.EXPECT().Build().Return(mockClient)
				mockClient.EXPECT().Endpoint().Return(mockEndpoint)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.refreshToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if !tt.wantErr {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockClientBuilder := httpclientv1.NewMockClientBuilder[*request.SDKLogin, *response.LoginResponse](ctrl)
				mockClient := httpclientv1.NewMockClient[*request.SDKLogin, *response.LoginResponse](ctrl)
				expectedToken := "new-test-token"
				expectedExpiresAt := time.Now().Unix() + 7200
				mockEndpoint := func(ctx context.Context, req httpclientv1.HTTPRequest[*request.SDKLogin]) (*response.LoginResponse, error) {
					return &response.LoginResponse{
						Token:     expectedToken,
						ExpiresAt: expectedExpiresAt,
					}, nil
				}
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().SDKLogin().Return(mockClientBuilder)
				mockClientBuilder.EXPECT().Build().Return(mockClient)
				mockClient.EXPECT().Endpoint().Return(mockEndpoint)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.refreshToken()
				assert.NoError(t, err)
				assert.Equal(t, expectedToken, authMw.GetToken())
				assert.Equal(t, expectedExpiresAt, authMw.GetExpiresAt())
				return
			}
		})
	}
}

func Test_applyAuthOptions(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
		token     string
		expiresAt int64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "apply auth options with valid token",
			fields: fields{
				apiKey:    "test-key",
				apiSecret: "test-secret",
				token:     "test-token",
				expiresAt: time.Now().Unix() + 7200,
			},
		},
		{
			name: "apply auth options with expiring token",
			fields: fields{
				apiKey:    "test-key",
				apiSecret: "test-secret",
				token:     "test-token",
				expiresAt: time.Now().Unix() + 1800,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClientBuilder := httpclientv1.NewMockClientBuilder[httpclientv1.NilBody, []v1.GroupInfo](ctrl)

			mockClientBuilder.EXPECT().Use(gomock.Any()).DoAndReturn(func(fn func(httpclientv1.Endpoint[httpclientv1.NilBody, []v1.GroupInfo]) httpclientv1.Endpoint[httpclientv1.NilBody, []v1.GroupInfo]) httpclientv1.ClientBuilder[httpclientv1.NilBody, []v1.GroupInfo] {
				var capturedToken string
				testEndpoint := httpclientv1.NewEndpoint[httpclientv1.NilBody, []v1.GroupInfo](func(ctx context.Context, req interface{}) (interface{}, error) {
					token, ok := ctx.Value(RequestTokenKey).(string)
					if ok {
						capturedToken = token
					}
					return []v1.GroupInfo{}, nil
				})
				wrappedEp := fn(testEndpoint)
				ctx := context.Background()
				req := httpclientv1.HTTPRequest[httpclientv1.NilBody]{}
				_, _ = wrappedEp(ctx, req)
				assert.Equal(t, tt.fields.token, capturedToken)
				return mockClientBuilder
			}).Times(1)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{
				apiKey:    tt.fields.apiKey,
				apiSecret: tt.fields.apiSecret,
				token:     tt.fields.token,
				expiresAt: tt.fields.expiresAt,
			}

			got := applyAuthOptions(authMw, mockClientBuilder)
			assert.NotNil(t, got)
		})
	}
}

func Test_authMiddleware_passiveRefreshToken(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
		token     string
		expiresAt int64
	}
	type args struct {
		newToken     string
		newExpiresAt string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantToken     string
		wantExpiresAt int64
	}{
		{
			name: "passive refresh token with new token",
			fields: fields{
				apiKey:    "test-key",
				apiSecret: "test-secret",
				token:     "old-token",
				expiresAt: 1000,
			},
			args: args{
				newToken:     "new-token",
				newExpiresAt: "2000",
			},
			wantToken:     "new-token",
			wantExpiresAt: 2000,
		},
		{
			name: "passive refresh token without new token",
			fields: fields{
				apiKey:    "test-key",
				apiSecret: "test-secret",
				token:     "old-token",
				expiresAt: 1000,
			},
			args: args{
				newToken:     "",
				newExpiresAt: "2000",
			},
			wantToken:     "old-token",
			wantExpiresAt: 1000,
		},
		{
			name: "passive refresh token with invalid expires at",
			fields: fields{
				apiKey:    "test-key",
				apiSecret: "test-secret",
				token:     "old-token",
				expiresAt: 1000,
			},
			args: args{
				newToken:     "new-token",
				newExpiresAt: "invalid",
			},
			wantToken:     "old-token",
			wantExpiresAt: 1000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)

			authMw := &authMiddleware{
				txp:       mockFactory,
				apiKey:    tt.fields.apiKey,
				apiSecret: tt.fields.apiSecret,
				token:     tt.fields.token,
				expiresAt: tt.fields.expiresAt,
			}

			resp := httptest.NewRecorder()
			resp.Header().Set("new-token", tt.args.newToken)
			resp.Header().Set("new-expires-at", tt.args.newExpiresAt)

			ctx := context.Background()
			newCtx := authMw.passiveRefreshToken(ctx, resp.Result())

			authMw.mu.RLock()
			gotToken := authMw.token
			gotExpiresAt := authMw.expiresAt
			authMw.mu.RUnlock()

			assert.Equal(t, tt.wantToken, gotToken)
			assert.Equal(t, tt.wantExpiresAt, gotExpiresAt)

			if tt.args.newToken != "" && tt.args.newExpiresAt != "invalid" {
				tokenFromCtx, ok := newCtx.Value(RequestTokenKey).(string)
				assert.True(t, ok)
				assert.Equal(t, tt.wantToken, tokenFromCtx)
			}
		})
	}
}
