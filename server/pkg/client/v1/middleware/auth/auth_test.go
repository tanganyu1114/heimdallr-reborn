package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/utils"

	"github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/transport"

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
			errMsg:  "sdkChallenge client builder is nil",
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
			errMsg:  "sdkChallenge client builder is nil",
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
				mockSysUserTransport.EXPECT().GetSDKChallenge().Return(nil).AnyTimes()
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

func TestAuthMiddleware_HandleTokenRefreshFromResponse(t *testing.T) {
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

			newCtx := authMw.handleTokenRefreshFromResponse(tt.args.ctx, tt.args.response)

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

func Test_authMiddleware_buildClientOptions(t *testing.T) {
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
			wantLen: 2,
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

			opts := authMw.buildClientOptions()
			assert.Len(t, opts, tt.wantLen)
		})
	}
}

func Test_authMiddleware_injectTokenToContext(t *testing.T) {
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
				mockSysUserTransport.EXPECT().GetSDKChallenge().Return(nil).AnyTimes()
			}

			mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
			authMw := mw(mockFactory).(*authMiddleware)

			authMw.mu.Lock()
			authMw.token = tt.fields.token
			authMw.expiresAt = tt.fields.expiresAt
			authMw.mu.Unlock()

			got, err := authMw.injectTokenToContext(tt.args.ctx)
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

func Test_authMiddleware_fetchNewToken(t *testing.T) {
	// Generate real RSA keys for testing
	utils.GenerateRSAKeys()
	testPublicKey, _, _ := utils.GetPublicKeyWithChallenge("test-api-key")

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
			name: "sdkChallenge client builder is nil",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
			},
			wantErr: true,
			errMsg:  "sdkChallenge client builder is nil",
		},
		{
			name: "sdkChallenge client is nil",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
			},
			wantErr: true,
			errMsg:  "sdkChallenge client is nil",
		},
		{
			name: "sdk challenge endpoint is nil",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
			},
			wantErr: true,
			errMsg:  "sdk challenge endpoint is nil",
		},
		{
			name: "sdkChallenge request fails",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
			},
			wantErr: true,
			errMsg:  "challenge request failed",
		},
		{
			name: "sdkChallenge response parse fails",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
			},
			wantErr: true,
			errMsg:  "invalid character",
		},
		{
			name: "failed to encrypt SDK login request",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
			},
			wantErr: true,
			errMsg:  "failed to parse PEM block containing the public key",
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
			name: "login response code is non-zero",
			fields: fields{
				apiKey:    "key",
				apiSecret: "secret",
			},
			wantErr: true,
			errMsg:  "internal server error",
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

				err := authMw.fetchNewToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "sysUser transport is nil" {
				mockFactory.EXPECT().SysUsers().Return(nil)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.fetchNewToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "sdkChallenge client builder is nil" {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().GetSDKChallenge().Return(nil)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.fetchNewToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "sdkChallenge client is nil" {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockChallengeBuilder := httpclientv1.NewMockClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().GetSDKChallenge().Return(mockChallengeBuilder)
				mockChallengeBuilder.EXPECT().Build().Return(nil)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.fetchNewToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "sdk challenge endpoint is nil" {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockChallengeBuilder := httpclientv1.NewMockClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeClient := httpclientv1.NewMockClient[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().GetSDKChallenge().Return(mockChallengeBuilder)
				mockChallengeBuilder.EXPECT().Build().Return(mockChallengeClient)
				mockChallengeClient.EXPECT().Endpoint().Return(nil)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.fetchNewToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "challenge request failed" {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockChallengeBuilder := httpclientv1.NewMockClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeClient := httpclientv1.NewMockClient[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeEndpoint := func(ctx context.Context, req httpclientv1.HTTPRequest[*request.SDKChallengeRequest]) (modelclientv1.ResponseBody[*response.SDKChallengeResponse], error) {
					return modelclientv1.ResponseBody[*response.SDKChallengeResponse]{}, errors.New("challenge request failed")
				}
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().GetSDKChallenge().Return(mockChallengeBuilder)
				mockChallengeBuilder.EXPECT().Build().Return(mockChallengeClient)
				mockChallengeClient.EXPECT().Endpoint().Return(mockChallengeEndpoint)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.fetchNewToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "invalid character" {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockChallengeBuilder := httpclientv1.NewMockClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeClient := httpclientv1.NewMockClient[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeEndpoint := func(ctx context.Context, req httpclientv1.HTTPRequest[*request.SDKChallengeRequest]) (modelclientv1.ResponseBody[*response.SDKChallengeResponse], error) {
					return modelclientv1.ResponseBody[*response.SDKChallengeResponse]{
						Data: json.RawMessage("invalid-json"),
					}, nil
				}
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().GetSDKChallenge().Return(mockChallengeBuilder)
				mockChallengeBuilder.EXPECT().Build().Return(mockChallengeClient)
				mockChallengeClient.EXPECT().Endpoint().Return(mockChallengeEndpoint)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.fetchNewToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "failed to parse PEM block containing the public key" {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockChallengeBuilder := httpclientv1.NewMockClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeClient := httpclientv1.NewMockClient[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeEndpoint := func(ctx context.Context, req httpclientv1.HTTPRequest[*request.SDKChallengeRequest]) (modelclientv1.ResponseBody[*response.SDKChallengeResponse], error) {
					return modelclientv1.ResponseBody[*response.SDKChallengeResponse]{
						Data: json.RawMessage(`{"public_key":"invalid-key","challenge":"test-challenge"}`),
					}, nil
				}
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().GetSDKChallenge().Return(mockChallengeBuilder)
				mockChallengeBuilder.EXPECT().Build().Return(mockChallengeClient)
				mockChallengeClient.EXPECT().Endpoint().Return(mockChallengeEndpoint)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.fetchNewToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "sdkLogin client builder is nil" {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockChallengeBuilder := httpclientv1.NewMockClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeClient := httpclientv1.NewMockClient[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeEndpoint := func(ctx context.Context, req httpclientv1.HTTPRequest[*request.SDKChallengeRequest]) (modelclientv1.ResponseBody[*response.SDKChallengeResponse], error) {
					challengeData := &response.SDKChallengeResponse{
						PublicKey: testPublicKey,
						Challenge: "test-challenge",
					}
					rawData, _ := json.Marshal(challengeData)
					return modelclientv1.ResponseBody[*response.SDKChallengeResponse]{
						Data: rawData,
					}, nil
				}
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().GetSDKChallenge().Return(mockChallengeBuilder)
				mockChallengeBuilder.EXPECT().Build().Return(mockChallengeClient)
				mockChallengeClient.EXPECT().Endpoint().Return(mockChallengeEndpoint)
				mockSysUserTransport.EXPECT().SDKLogin().Return(nil)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.fetchNewToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "sdkLogin client is nil" {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockChallengeBuilder := httpclientv1.NewMockClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeClient := httpclientv1.NewMockClient[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeEndpoint := func(ctx context.Context, req httpclientv1.HTTPRequest[*request.SDKChallengeRequest]) (modelclientv1.ResponseBody[*response.SDKChallengeResponse], error) {
					challengeData, _ := json.Marshal(&response.SDKChallengeResponse{
						PublicKey: testPublicKey,
						Challenge: "test-challenge",
					})
					return modelclientv1.ResponseBody[*response.SDKChallengeResponse]{
						Data: challengeData,
					}, nil
				}
				mockLoginBuilder := httpclientv1.NewMockClientBuilder[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]](ctrl)
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().GetSDKChallenge().Return(mockChallengeBuilder)
				mockChallengeBuilder.EXPECT().Build().Return(mockChallengeClient)
				mockChallengeClient.EXPECT().Endpoint().Return(mockChallengeEndpoint)
				mockSysUserTransport.EXPECT().SDKLogin().Return(mockLoginBuilder)
				mockLoginBuilder.EXPECT().Build().Return(nil)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.fetchNewToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "sdk login endpoint is nil" {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockChallengeBuilder := httpclientv1.NewMockClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeClient := httpclientv1.NewMockClient[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeEndpoint := func(ctx context.Context, req httpclientv1.HTTPRequest[*request.SDKChallengeRequest]) (modelclientv1.ResponseBody[*response.SDKChallengeResponse], error) {
					challengeData, _ := json.Marshal(&response.SDKChallengeResponse{
						PublicKey: testPublicKey,
						Challenge: "test-challenge",
					})
					return modelclientv1.ResponseBody[*response.SDKChallengeResponse]{
						Data: challengeData,
					}, nil
				}
				mockLoginBuilder := httpclientv1.NewMockClientBuilder[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]](ctrl)
				mockLoginClient := httpclientv1.NewMockClient[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]](ctrl)
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().GetSDKChallenge().Return(mockChallengeBuilder)
				mockChallengeBuilder.EXPECT().Build().Return(mockChallengeClient)
				mockChallengeClient.EXPECT().Endpoint().Return(mockChallengeEndpoint)
				mockSysUserTransport.EXPECT().SDKLogin().Return(mockLoginBuilder)
				mockLoginBuilder.EXPECT().Build().Return(mockLoginClient)
				mockLoginClient.EXPECT().Endpoint().Return(nil)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.fetchNewToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.errMsg == "login error" {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockChallengeBuilder := httpclientv1.NewMockClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeClient := httpclientv1.NewMockClient[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeEndpoint := func(ctx context.Context, req httpclientv1.HTTPRequest[*request.SDKChallengeRequest]) (modelclientv1.ResponseBody[*response.SDKChallengeResponse], error) {
					challengeData, _ := json.Marshal(&response.SDKChallengeResponse{
						PublicKey: testPublicKey,
						Challenge: "test-challenge",
					})
					return modelclientv1.ResponseBody[*response.SDKChallengeResponse]{
						Data: challengeData,
					}, nil
				}
				mockLoginBuilder := httpclientv1.NewMockClientBuilder[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]](ctrl)
				mockLoginClient := httpclientv1.NewMockClient[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]](ctrl)
				mockLoginEndpoint := func(ctx context.Context, req httpclientv1.HTTPRequest[*request.EncryptedLoginRequest]) (modelclientv1.ResponseBody[*response.LoginResponse], error) {
					return modelclientv1.ResponseBody[*response.LoginResponse]{}, errors.New("login error")
				}
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().GetSDKChallenge().Return(mockChallengeBuilder)
				mockChallengeBuilder.EXPECT().Build().Return(mockChallengeClient)
				mockChallengeClient.EXPECT().Endpoint().Return(mockChallengeEndpoint)
				mockSysUserTransport.EXPECT().SDKLogin().Return(mockLoginBuilder)
				mockLoginBuilder.EXPECT().Build().Return(mockLoginClient)
				mockLoginClient.EXPECT().Endpoint().Return(mockLoginEndpoint)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.fetchNewToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if tt.name == "login response code is non-zero" {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockChallengeBuilder := httpclientv1.NewMockClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeClient := httpclientv1.NewMockClient[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeEndpoint := func(ctx context.Context, req httpclientv1.HTTPRequest[*request.SDKChallengeRequest]) (modelclientv1.ResponseBody[*response.SDKChallengeResponse], error) {
					challengeData, _ := json.Marshal(&response.SDKChallengeResponse{
						PublicKey: testPublicKey,
						Challenge: "test-challenge",
					})
					return modelclientv1.ResponseBody[*response.SDKChallengeResponse]{
						Data: challengeData,
					}, nil
				}
				mockLoginBuilder := httpclientv1.NewMockClientBuilder[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]](ctrl)
				mockLoginClient := httpclientv1.NewMockClient[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]](ctrl)
				mockLoginEndpoint := httpclientv1.Endpoint[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]](func(ctx context.Context, req httpclientv1.HTTPRequest[*request.EncryptedLoginRequest]) (resp modelclientv1.ResponseBody[*response.LoginResponse], err error) {
					return modelclientv1.ResponseBody[*response.LoginResponse]{
						Code:    500,
						Message: "internal server error",
					}, nil
				})
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().GetSDKChallenge().Return(mockChallengeBuilder)
				mockChallengeBuilder.EXPECT().Build().Return(mockChallengeClient)
				mockChallengeClient.EXPECT().Endpoint().Return(mockChallengeEndpoint)
				mockSysUserTransport.EXPECT().SDKLogin().Return(mockLoginBuilder)
				mockLoginBuilder.EXPECT().Build().Return(mockLoginClient)
				mockLoginClient.EXPECT().Endpoint().Return(mockLoginEndpoint)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.fetchNewToken()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			if !tt.wantErr {
				mockSysUserTransport := transport.NewMockSysUserTransport(ctrl)
				mockChallengeBuilder := httpclientv1.NewMockClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeClient := httpclientv1.NewMockClient[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)
				mockChallengeEndpoint := func(ctx context.Context, req httpclientv1.HTTPRequest[*request.SDKChallengeRequest]) (modelclientv1.ResponseBody[*response.SDKChallengeResponse], error) {
					challengeData, _ := json.Marshal(&response.SDKChallengeResponse{
						PublicKey: testPublicKey,
						Challenge: "test-challenge",
					})
					return modelclientv1.ResponseBody[*response.SDKChallengeResponse]{
						Data: challengeData,
					}, nil
				}
				mockLoginBuilder := httpclientv1.NewMockClientBuilder[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]](ctrl)
				mockLoginClient := httpclientv1.NewMockClient[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]](ctrl)
				expectedToken := "new-test-token"
				expectedExpiresAt := time.Now().Unix() + 7200
				mockLoginEndpoint := func(ctx context.Context, req httpclientv1.HTTPRequest[*request.EncryptedLoginRequest]) (modelclientv1.ResponseBody[*response.LoginResponse], error) {
					data, _ := json.Marshal(&response.LoginResponse{
						Token:     expectedToken,
						ExpiresAt: expectedExpiresAt,
					})
					return modelclientv1.ResponseBody[*response.LoginResponse]{
						Data: data,
					}, nil
				}
				mockFactory.EXPECT().SysUsers().Return(mockSysUserTransport)
				mockSysUserTransport.EXPECT().GetSDKChallenge().Return(mockChallengeBuilder)
				mockChallengeBuilder.EXPECT().Build().Return(mockChallengeClient)
				mockChallengeClient.EXPECT().Endpoint().Return(mockChallengeEndpoint)
				mockSysUserTransport.EXPECT().SDKLogin().Return(mockLoginBuilder)
				mockLoginBuilder.EXPECT().Build().Return(mockLoginClient)
				mockLoginClient.EXPECT().Endpoint().Return(mockLoginEndpoint)

				mw := AuthMiddleware(tt.fields.apiKey, tt.fields.apiSecret, tt.fields.bufferTime)
				authMw := mw(mockFactory).(*authMiddleware)

				err := authMw.fetchNewToken()
				assert.NoError(t, err)
				assert.Equal(t, expectedToken, authMw.GetToken())
				assert.Equal(t, expectedExpiresAt, authMw.GetExpiresAt())
				return
			}
		})
	}
}

func Test_wrapWithAuthOptions(t *testing.T) {
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

			mockClientBuilder := httpclientv1.NewMockClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.GroupInfo]](ctrl)

			mockClientBuilder.EXPECT().Use(gomock.Any()).DoAndReturn(func(fn func(httpclientv1.Endpoint[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.GroupInfo]]) httpclientv1.Endpoint[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.GroupInfo]]) httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.GroupInfo]] {
				var capturedToken string
				testEndpoint := httpclientv1.NewEndpoint[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.GroupInfo]](func(ctx context.Context, request interface{}) (response interface{}, err error) {
					token, ok := ctx.Value(RequestTokenKey).(string)
					if ok {
						capturedToken = token
					}
					data, _ := json.Marshal([]v1.GroupInfo{})
					return modelclientv1.ResponseBody[[]v1.GroupInfo]{Data: data}, nil
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

			got := wrapWithAuth(authMw, mockClientBuilder)
			assert.NotNil(t, got)
		})
	}
}

func Test_authMiddleware_handleTokenRefreshFromResponse(t *testing.T) {
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
			newCtx := authMw.handleTokenRefreshFromResponse(ctx, resp.Result())

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

func TestAuthMiddleware_FullRequestResponseFlow(t *testing.T) {
	tests := []struct {
		name           string
		initialToken   string
		initialExpires int64
		newToken       string
		newExpiresAt   string
		wantReqToken   string
		wantNewToken   string
		wantNewExpires int64
	}{
		{
			name:           "token injected into request header and new token captured from response",
			initialToken:   "initial-token",
			initialExpires: time.Now().Unix() + 7200,
			newToken:       "refreshed-token",
			newExpiresAt:   strconv.FormatInt(time.Now().Unix()+14400, 10),
			wantReqToken:   "initial-token",
			wantNewToken:   "refreshed-token",
			wantNewExpires: time.Now().Unix() + 14400,
		},
		{
			name:           "token injected but no new token in response",
			initialToken:   "initial-token",
			initialExpires: time.Now().Unix() + 7200,
			newToken:       "",
			newExpiresAt:   "",
			wantReqToken:   "initial-token",
			wantNewToken:   "initial-token",
			wantNewExpires: time.Now().Unix() + 7200,
		},
		{
			name:           "token injected but older token in response ignored",
			initialToken:   "initial-token",
			initialExpires: time.Now().Unix() + 7200,
			newToken:       "older-token",
			newExpiresAt:   strconv.FormatInt(time.Now().Unix()+3600, 10),
			wantReqToken:   "initial-token",
			wantNewToken:   "initial-token",
			wantNewExpires: time.Now().Unix() + 7200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedRequestToken string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedRequestToken = r.Header.Get(RequestTokenKey)
				w.Header().Set("Content-Type", "application/json")
				if tt.newToken != "" {
					w.Header().Set("new-token", tt.newToken)
					w.Header().Set("new-expires-at", tt.newExpiresAt)
				}
				data, _ := json.Marshal(modelclientv1.ResponseBody[[]v1.GroupInfo]{
					Data: json.RawMessage(`[]`),
				})
				w.Write(data)
			}))
			defer server.Close()

			authMw := &authMiddleware{
				apiKey:     "test-key",
				apiSecret:  "test-secret",
				bufferTime: 3600,
				token:      tt.initialToken,
				expiresAt:  tt.initialExpires,
			}

			realTxp, err := transport.NewTransport(server.Client(), server.URL)
			assert.NoError(t, err)
			authMw.txp = realTxp

			// Get the auth-wrapped ClientBuilder directly from the middleware
			clientBuilder := newAgentInfoMiddleware(authMw).Get()
			client := clientBuilder.Build()

			resp, err := client.Endpoint()(context.Background(), httpclientv1.HTTPRequest[httpclientv1.NilBody]{})
			assert.NoError(t, err)
			assert.NotNil(t, resp)

			assert.Equal(t, tt.wantReqToken, capturedRequestToken, "HTTP request header should contain x-token injected by auth middleware")

			if tt.newToken != "" {
				newExpires, _ := strconv.ParseInt(tt.newExpiresAt, 10, 64)
				if newExpires > tt.initialExpires {
					assert.Equal(t, tt.wantNewToken, authMw.GetToken(), "auth middleware should update token from response header")
					assert.Equal(t, newExpires, authMw.GetExpiresAt(), "auth middleware should update expiresAt from response header")
				} else {
					assert.Equal(t, tt.wantNewToken, authMw.GetToken(), "auth middleware should keep original token when new expiresAt is older")
					assert.Equal(t, tt.initialExpires, authMw.GetExpiresAt(), "auth middleware should keep original expiresAt when new one is older")
				}
			} else {
				assert.Equal(t, tt.wantNewToken, authMw.GetToken(), "auth middleware should keep original token when no new token in response")
				assert.Equal(t, tt.initialExpires, authMw.GetExpiresAt(), "auth middleware should keep original expiresAt when no new token in response")
			}
		})
	}
}
