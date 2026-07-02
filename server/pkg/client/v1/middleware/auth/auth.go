package auth

import (
	"context"
	"gin-vue-admin/pkg/client/v1/middleware"
	txpclientv1 "gin-vue-admin/pkg/client/v1/transport"
	"net/http"
	"strconv"
	"sync"
	"time"

	"gin-vue-admin/model/request"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	http_transport "github.com/go-kit/kit/transport/http"
	"github.com/marmotedu/errors"
)

const (
	RequestTokenKey = "X-Token"
	MiddlewareName  = "AuthMiddleware"
)

func withRequestToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, RequestTokenKey, token)
}

func AuthMiddleware(apiKey, apiSecret string, bufferTime int64) middleware.Middleware {
	return func(factory txpclientv1.Factory) txpclientv1.Factory {
		return &authMiddleware{
			txp:        factory,
			apiKey:     apiKey,
			apiSecret:  apiSecret,
			bufferTime: bufferTime,
		}
	}
}

type authMiddleware struct {
	txp        txpclientv1.Factory
	apiKey     string
	apiSecret  string
	bufferTime int64
	mu         sync.RWMutex
	token      string
	expiresAt  int64
}

func (a *authMiddleware) SysUsers() txpclientv1.SysUserTransport {
	return a.txp.SysUsers()
}

func (a *authMiddleware) AgentInfos() txpclientv1.AgentInfoTransport {
	return newAgentInfoMiddleware(a)
}

func (a *authMiddleware) Groups() txpclientv1.GroupTransport {
	return newGroupMiddleware(a)
}

func (a *authMiddleware) Hosts() txpclientv1.HostTransport {
	return newHostMiddleware(a)
}

func (a *authMiddleware) WebServerConfigs() txpclientv1.WebServerConfigTransport {
	return newWebServerConfigMiddleware(a)
}

func (a *authMiddleware) WebServerBinCMDs() txpclientv1.WebServerBinCMDTransport {
	return newWebServerBinCMDMiddleware(a)
}

func (a *authMiddleware) WebServerStatistics() txpclientv1.WebServerStatisticsTransport {
	return newWebServerStatisticsMiddleware(a)
}

func (a *authMiddleware) authOptions() []http_transport.ClientOption {
	return []http_transport.ClientOption{
		http_transport.ClientAfter(a.passiveRefreshToken),
	}
}

func (a *authMiddleware) ensureValidTokenToCtx(ctx context.Context) (context.Context, error) {
	if err := a.ensureValidToken(); err != nil {
		logV1.Errorf("failed to ensure valid token: %v", err)
		return ctx, err
	}
	a.mu.RLock()
	token := a.token
	a.mu.RUnlock()
	return withRequestToken(ctx, token), nil
}

func (a *authMiddleware) passiveRefreshToken(ctx context.Context, response *http.Response) context.Context {
	newToken := response.Header.Get("new-token")
	if newToken == "" {
		return ctx
	}

	newExpiresAtStr := response.Header.Get("new-expires-at")
	newExpiresAt, err := strconv.ParseInt(newExpiresAtStr, 10, 64)
	if err != nil {
		logV1.Errorf("failed to parse new-expires-at header: %v", err)
		return ctx
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if newExpiresAt > a.expiresAt {
		a.token = newToken
		a.expiresAt = newExpiresAt
		ctx = withRequestToken(ctx, newToken)
	}

	return ctx
}

func applyAuthOptions[REQ any, RESP any](mw *authMiddleware, clientBuilder httpclientv1.ClientBuilder[REQ, RESP]) httpclientv1.ClientBuilder[REQ, RESP] {
	return clientBuilder.
		Use(func(ep httpclientv1.Endpoint[REQ, RESP]) httpclientv1.Endpoint[REQ, RESP] {
			return func(ctx context.Context, request httpclientv1.HTTPRequest[REQ]) (response RESP, err error) {
				ctx, err = mw.ensureValidTokenToCtx(ctx)
				if err != nil {
					return
				}
				return ep(ctx, request)
			}
		}).
		WithOptions(mw.authOptions()...)
}

func (a *authMiddleware) refreshToken() error {
	a.mu.RLock()
	currentAPIKey := a.apiKey
	currentAPISecret := a.apiSecret
	txp := a.txp
	a.mu.RUnlock()

	if txp == nil {
		return errors.New("transport is nil")
	}

	sysUsers := txp.SysUsers()
	if sysUsers == nil {
		return errors.New("sysUser transport is nil")
	}

	sdkLoginBuilder := sysUsers.SDKLogin()
	if sdkLoginBuilder == nil {
		return errors.New("sdkLogin client builder is nil")
	}

	client := sdkLoginBuilder.Build()
	if client == nil {
		return errors.New("sdkLogin client is nil")
	}

	sdkLoginEP := client.Endpoint()
	if sdkLoginEP == nil {
		return errors.New("sdk login endpoint is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	loginReq := &request.SDKLogin{
		APIKey:    currentAPIKey,
		APISecret: currentAPISecret,
	}

	req := httpclientv1.HTTPRequest[*request.SDKLogin]{
		Body: loginReq,
	}

	loginResp, err := sdkLoginEP(ctx, req)
	if err != nil {
		logV1.Errorf("logining and refreshing token failed: %v", err)
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.token = loginResp.Token
	a.expiresAt = loginResp.ExpiresAt

	return nil
}

func (a *authMiddleware) ensureValidToken() error {
	a.mu.RLock()
	expiresAt := a.expiresAt
	bufferTime := a.bufferTime
	tokenExists := a.token != ""
	a.mu.RUnlock()

	if !tokenExists || expiresAt-time.Now().Unix() < bufferTime {
		if err := a.refreshToken(); err != nil {
			logV1.Errorf("failed to ensure valid token: %v", err)
			return err
		}
	}
	return nil
}

func (a *authMiddleware) GetToken() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.token
}

func (a *authMiddleware) GetExpiresAt() int64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.expiresAt
}

func (a *authMiddleware) GetAPIKey() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.apiKey
}

func (a *authMiddleware) GetAPISecret() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.apiSecret
}
