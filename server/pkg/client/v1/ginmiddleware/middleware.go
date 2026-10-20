package ginmiddleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1"
)

const (
	ClaimsKey          = "claims"
	SDKLoginKey        = "is_sdk_login"
	TokenHeader        = "x-token"
	UserIDHeader       = "x-user-id"
	NewTokenHeader     = "new-token"
	NewExpiresAtHeader = "new-expires-at"
)

type GinMiddleware struct {
	client   *v1.Client
	apiGroup string
}

func New(client *v1.Client, apiGroup string) (*GinMiddleware, error) {
	if client == nil {
		return nil, ErrClientNil
	}
	if apiGroup == "" {
		return nil, ErrApiGroupEmpty
	}
	return &GinMiddleware{client: client, apiGroup: apiGroup}, nil
}

func (gm *GinMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(TokenHeader)
		if token == "" {
			respondError(c, gin.H{"reload": true}, "未登录或非法访问")
			c.Abort()
			return
		}

		// Call remote SDK service to validate token
		resp, err := gm.client.Authz().ValidateToken(token)
		if err != nil {
			handleTokenError(c, err)
			return
		}

		claims := resp.Claims
		if gm.shouldRefreshToken(claims) {
			gm.refreshToken(c, token, claims)
		}

		c.Set(ClaimsKey, claims)
		c.Next()
	}
}

func (gm *GinMiddleware) CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get(ClaimsKey)
		if !exists {
			respondError(c, gin.H{}, "未认证")
			c.Abort()
			return
		}

		customClaims := claims.(*request.CustomClaims)
		obj := c.Request.URL.RequestURI()
		act := c.Request.Method
		sub := customClaims.AuthorityId

		// Call remote SDK service to check permission
		resp, err := gm.client.Authz().CheckPermission(sub, obj, act)
		if err != nil {
			respondError(c, gin.H{}, "权限检查失败")
			c.Abort()
			return
		}

		if !resp.Allowed {
			respondError(c, gin.H{}, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (gm *GinMiddleware) OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		var userID int

		if c.Request.Method != http.MethodGet {
			body = readBody(c)
		}

		userID = extractUserID(c)

		record := &request.LogOperationRequest{
			IP:     c.ClientIP(),
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Agent:  c.Request.UserAgent(),
			Body:   string(body),
			UserID: userID,
		}

		if isSDKLogin, ok := c.Get(SDKLoginKey); ok && isSDKLogin.(bool) {
			record.IsSDKLogin = true
		}

		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		start := time.Now()
		c.Next()

		record.Latency = time.Since(start).Nanoseconds()
		record.Status = writer.Status()
		record.ErrorMessage = extractError(c)

		contentType := writer.Header().Get("Content-Type")
		if isBinaryResponse(contentType) {
			record.Resp = "[Binary File - Not Logged]"
		} else {
			record.Resp = writer.body.String()
		}

		go gm.logOperation(record)
	}
}

func (gm *GinMiddleware) logOperation(record *request.LogOperationRequest) {
	if err := gm.client.Authz().LogOperation(record); err != nil {
		// Log error but don't fail the request
		logV1.Warnf("logOperation err: %+v", err)
	}
}

func respondError(c *gin.Context, data gin.H, msg string) {
	c.JSON(401, gin.H{
		"success": false,
		"data":    data,
		"message": msg,
	})
}

func handleTokenError(c *gin.Context, err error) {
	reload := gin.H{"reload": true}
	if strings.Contains(err.Error(), "expired") {
		respondError(c, reload, "授权已过期")
	} else {
		respondError(c, reload, err.Error())
	}
	c.Abort()
}

func (gm *GinMiddleware) shouldRefreshToken(claims *request.CustomClaims) bool {
	return claims.ExpiresAt-time.Now().Unix() < claims.BufferTime
}

func (gm *GinMiddleware) refreshToken(c *gin.Context, token string, claims *request.CustomClaims) {
	// Calculate new expiration time (7 days from now, same as original jwt.go)
	newExpiresAt := time.Now().Unix() + 60*60*24*7

	// Call remote SDK service to refresh token
	newTokenResp, err := gm.client.Authz().RefreshToken(token, newExpiresAt)
	if err != nil {
		logV1.Errorf("refreshToken err: %+v", err)
		return
	}

	// Set new token and expires-at in response headers
	c.Header(NewTokenHeader, newTokenResp.Token)
	c.Header(NewExpiresAtHeader, strconv.FormatInt(newTokenResp.ExpiresAt, 10))
}

func readBody(c *gin.Context) []byte {
	body, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body
}

func extractUserID(c *gin.Context) int {
	if claims, ok := c.Get(ClaimsKey); ok {
		customClaims := claims.(*request.CustomClaims)
		return int(customClaims.ID)
	}

	idStr := c.Request.Header.Get(UserIDHeader)
	id, _ := strconv.Atoi(idStr)
	return id
}

func extractError(c *gin.Context) string {
	return c.Errors.ByType(gin.ErrorTypePrivate).String()
}

func isBinaryResponse(contentType string) bool {
	binaryTypes := []string{
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/octet-stream",
		"application/pdf",
	}

	for _, bt := range binaryTypes {
		if contentType == bt || strings.HasPrefix(contentType, "image/") {
			return true
		}
	}
	return false
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
