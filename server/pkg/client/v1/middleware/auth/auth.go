package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/middleware"
	txpclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/transport"
	"github.com/tanganyu1114/heimdallr-reborn/server/utils"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	http_transport "github.com/go-kit/kit/transport/http"
	"github.com/marmotedu/errors"
)

const (
	RequestTokenKey = "X-Token"
	MiddlewareName  = "AuthMiddleware"
)

type authMiddleware struct {
	txp        txpclientv1.Factory
	apiKey     string
	apiSecret  string
	bufferTime int64
	mu         sync.RWMutex
	token      string
	expiresAt  int64
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

func (a *authMiddleware) fetchNewToken() error {
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

	// Step 1: Get SDK challenge (public key + challenge)
	sdkChallengeBuilder := sysUsers.GetSDKChallenge()
	if sdkChallengeBuilder == nil {
		return errors.New("sdkChallenge client builder is nil")
	}

	sdkChallengeClient := sdkChallengeBuilder.Build()
	if sdkChallengeClient == nil {
		return errors.New("sdkChallenge client is nil")
	}

	sdkChallengeEP := sdkChallengeClient.Endpoint()
	if sdkChallengeEP == nil {
		return errors.New("sdk challenge endpoint is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	challengeReq := &request.SDKChallengeRequest{
		APIKey: currentAPIKey,
	}

	challengeHTTPReq := httpclientv1.HTTPRequest[*request.SDKChallengeRequest]{
		Body: challengeReq,
	}

	challengeResp, err := sdkChallengeEP(ctx, challengeHTTPReq)
	if err != nil {
		logV1.Errorf("failed to get SDK challenge: %v", err)
		return err
	}

	challengeRespData, err := challengeResp.Response()
	if err != nil {
		logV1.Errorf("failed to parse SDK challenge response: %v", err)
		return err
	}

	// Step 2: Build SDK login request with challenge and encrypt it
	loginReq := &request.SDKLogin{
		APIKey:    currentAPIKey,
		APISecret: currentAPISecret,
		Challenge: challengeRespData.Challenge,
	}

	// Serialize login request to JSON
	loginJSON, err := json.Marshal(loginReq)
	if err != nil {
		logV1.Errorf("failed to marshal SDK login request: %v", err)
		return err
	}

	// Encrypt the login request with the public key
	encryptedData, err := utils.RSAEncrypt(challengeRespData.PublicKey, string(loginJSON))
	if err != nil {
		logV1.Errorf("failed to encrypt SDK login request: %v", err)
		return err
	}

	// Step 3: Send encrypted SDK login request
	sdkLoginBuilder := sysUsers.SDKLogin()
	if sdkLoginBuilder == nil {
		return errors.New("sdkLogin client builder is nil")
	}

	sdkLoginClient := sdkLoginBuilder.Build()
	if sdkLoginClient == nil {
		return errors.New("sdkLogin client is nil")
	}

	sdkLoginEP := sdkLoginClient.Endpoint()
	if sdkLoginEP == nil {
		return errors.New("sdk login endpoint is nil")
	}

	// Create encrypted request wrapper
	encryptedReq := request.EncryptedLoginRequest{
		EncryptedData: encryptedData,
	}

	req := httpclientv1.HTTPRequest[*request.EncryptedLoginRequest]{
		Body: &encryptedReq,
	}

	resp, err := sdkLoginEP(ctx, req)
	if err != nil {
		logV1.Errorf("failed to login and refresh token: %v", err)
		return err
	}

	loginResp, err := resp.Response()
	if err != nil {
		logV1.Errorf("failed to login and refresh token: %v", err)
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
		if err := a.fetchNewToken(); err != nil {
			logV1.Errorf("failed to ensure valid token: %v", err)
			return err
		}
	}
	return nil
}

func (a *authMiddleware) injectTokenToContext(ctx context.Context) (context.Context, error) {
	if err := a.ensureValidToken(); err != nil {
		logV1.Errorf("failed to inject token to context: %v", err)
		return ctx, err
	}
	a.mu.RLock()
	token := a.token
	a.mu.RUnlock()
	return context.WithValue(ctx, RequestTokenKey, token), nil
}

func (a *authMiddleware) buildClientOptions() []http_transport.ClientOption {
	return []http_transport.ClientOption{
		http_transport.ClientBefore(a.injectTokenToHeader),
		http_transport.ClientAfter(a.handleTokenRefreshFromResponse),
	}
}

func (a *authMiddleware) injectTokenToHeader(ctx context.Context, req *http.Request) context.Context {
	if token, ok := ctx.Value(RequestTokenKey).(string); ok {
		req.Header.Set(RequestTokenKey, token)
	}
	return ctx
}

func (a *authMiddleware) handleTokenRefreshFromResponse(ctx context.Context, response *http.Response) context.Context {
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
		ctx = context.WithValue(ctx, RequestTokenKey, newToken)
	}

	return ctx
}

func wrapWithAuth[REQ any, RESP any](mw *authMiddleware, clientBuilder httpclientv1.ClientBuilder[REQ, RESP]) httpclientv1.ClientBuilder[REQ, RESP] {
	return clientBuilder.
		Use(func(ep httpclientv1.Endpoint[REQ, RESP]) httpclientv1.Endpoint[REQ, RESP] {
			return func(ctx context.Context, request httpclientv1.HTTPRequest[REQ]) (response RESP, err error) {
				ctx, err = mw.injectTokenToContext(ctx)
				if err != nil {
					return
				}
				return ep(ctx, request)
			}
		}).
		WithOptions(mw.buildClientOptions()...)
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
