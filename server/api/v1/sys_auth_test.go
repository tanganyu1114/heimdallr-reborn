package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tanganyu1114/heimdallr-reborn/server/global"
	"github.com/tanganyu1114/heimdallr-reborn/server/middleware"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func init() {
	// Initialize global logger for tests
	global.GVA_LOG = zap.NewNop()
}

// Helper function to generate a valid JWT token for testing
func generateTestToken() string {
	j := middleware.NewJWT()
	claims := request.CustomClaims{
		UUID:        uuid.NewV4(),
		ID:          1,
		Username:    "testuser",
		NickName:    "Test User",
		AuthorityId: "admin",
		BufferTime:  86400,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,
			ExpiresAt: time.Now().Unix() + 86400,
			Issuer:    "test",
		},
	}
	token, _ := j.CreateToken(claims)
	return token
}

// TestValidateToken_Blacklist tests the blacklist checking logic in ValidateToken
func TestValidateToken_Blacklist(t *testing.T) {
	// Save original checker to restore later
	originalChecker := DefaultBlacklistChecker
	defer func() {
		DefaultBlacklistChecker = originalChecker
	}()

	tests := []struct {
		name          string
		isBlacklisted bool
		expectedCode  int
		expectedMsg   string
	}{
		{
			name:          "token not blacklisted",
			isBlacklisted: false,
			expectedCode:  response.SUCCESS,
			expectedMsg:   "",
		},
		{
			name:          "token is blacklisted",
			isBlacklisted: true,
			expectedCode:  response.ERROR,
			expectedMsg:   "token is blacklisted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock blacklist checker
			DefaultBlacklistChecker = func(jwt string) bool {
				return tt.isBlacklisted
			}

			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Generate a valid token for testing
			validToken := generateTestToken()
			requestBody := map[string]string{"token": validToken}
			body, _ := json.Marshal(requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/validateToken", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute
			ValidateToken(c)

			// Assert
			assert.Equal(t, http.StatusOK, w.Code)
			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.Code)
			if tt.expectedMsg != "" {
				assert.Equal(t, tt.expectedMsg, resp.Msg)
			}
		})
	}
}

// TestRefreshToken_Blacklist tests the blacklist checking logic in RefreshToken
func TestRefreshToken_Blacklist(t *testing.T) {
	// Save original checker to restore later
	originalChecker := DefaultBlacklistChecker
	defer func() {
		DefaultBlacklistChecker = originalChecker
	}()

	tests := []struct {
		name          string
		isBlacklisted bool
		expectedCode  int
		expectedMsg   string
	}{
		{
			name:          "token not blacklisted",
			isBlacklisted: false,
			expectedCode:  response.SUCCESS,
			expectedMsg:   "",
		},
		{
			name:          "token is blacklisted",
			isBlacklisted: true,
			expectedCode:  response.ERROR,
			expectedMsg:   "token is blacklisted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock blacklist checker
			DefaultBlacklistChecker = func(jwt string) bool {
				return tt.isBlacklisted
			}

			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Generate a valid token for testing
			validToken := generateTestToken()
			requestBody := map[string]interface{}{
				"token":     validToken,
				"expiresAt": time.Now().Unix() + 86400,
			}
			body, _ := json.Marshal(requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/refreshToken", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute
			RefreshToken(c)

			// Assert
			assert.Equal(t, http.StatusOK, w.Code)
			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.Code)
			if tt.expectedMsg != "" {
				assert.Equal(t, tt.expectedMsg, resp.Msg)
			}
		})
	}
}

// TestRefreshToken_CreateTokenError tests the token creation error path in RefreshToken
func TestRefreshToken_CreateTokenError(t *testing.T) {
	// Save original functions to restore later
	originalChecker := DefaultBlacklistChecker
	originalCreator := DefaultTokenCreator
	originalParser := DefaultTokenParser
	defer func() {
		DefaultBlacklistChecker = originalChecker
		DefaultTokenCreator = originalCreator
		DefaultTokenParser = originalParser
	}()

	// Setup mocks
	DefaultBlacklistChecker = func(jwt string) bool {
		return false // Token is not blacklisted
	}

	// Use real parser to generate a valid token
	j := middleware.NewJWT()
	DefaultTokenParser = j.ParseToken

	// Mock creator to return an error
	DefaultTokenCreator = func(claims request.CustomClaims) (string, error) {
		return "", assert.AnError // Return a test error
	}

	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Generate a valid token for testing
	validToken := generateTestToken()
	requestBody := map[string]interface{}{
		"token":     validToken,
		"expiresAt": time.Now().Unix() + 86400,
	}
	body, _ := json.Marshal(requestBody)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/refreshToken", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	RefreshToken(c)

	// Assert - HTTP status is always 200
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Should return error code and the error message
	assert.Equal(t, response.ERROR, resp.Code)
	assert.Contains(t, resp.Msg, "assert.AnError")
}

// TestRefreshToken_Success tests the successful token refresh path
func TestRefreshToken_Success(t *testing.T) {
	// Save original functions to restore later
	originalChecker := DefaultBlacklistChecker
	originalCreator := DefaultTokenCreator
	originalParser := DefaultTokenParser
	defer func() {
		DefaultBlacklistChecker = originalChecker
		DefaultTokenCreator = originalCreator
		DefaultTokenParser = originalParser
	}()

	// Setup mocks
	DefaultBlacklistChecker = func(jwt string) bool {
		return false // Token is not blacklisted
	}

	// Use real parser and creator
	j := middleware.NewJWT()
	DefaultTokenParser = j.ParseToken
	DefaultTokenCreator = j.CreateToken

	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Generate a valid token for testing
	validToken := generateTestToken()
	requestBody := map[string]interface{}{
		"token":     validToken,
		"expiresAt": time.Now().Unix() + 86400,
	}
	body, _ := json.Marshal(requestBody)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/refreshToken", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	RefreshToken(c)

	// Assert - HTTP status is always 200
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Should return success with new token
	assert.Equal(t, response.SUCCESS, resp.Code)
	assert.NotNil(t, resp.Data)

	// Verify the response contains a token
	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, data, "token")
	assert.NotEmpty(t, data["token"])
}

// TestValidateToken_JSONBindingError tests JSON binding error path in ValidateToken
func TestValidateToken_JSONBindingError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Send invalid JSON to trigger binding error
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/validateToken", bytes.NewReader([]byte("invalid-json")))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	ValidateToken(c)

	// Assert - HTTP status is always 200
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Should return error code due to JSON binding failure
	assert.Equal(t, response.ERROR, resp.Code)
	assert.NotEmpty(t, resp.Msg)
}

// TestRefreshToken_JSONBindingError tests JSON binding error path in RefreshToken
func TestRefreshToken_JSONBindingError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Send invalid JSON to trigger binding error
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/refreshToken", bytes.NewReader([]byte("invalid-json")))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	RefreshToken(c)

	// Assert - HTTP status is always 200
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Should return error code due to JSON binding failure
	assert.Equal(t, response.ERROR, resp.Code)
	assert.NotEmpty(t, resp.Msg)
}

// TestCheckPermission_JSONBindingError tests JSON binding error path in CheckPermission
func TestCheckPermission_JSONBindingError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Send invalid JSON to trigger binding error
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/checkPermission", bytes.NewReader([]byte("invalid-json")))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	CheckPermission(c)

	// Assert - HTTP status is always 200
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Should return error code due to JSON binding failure
	assert.Equal(t, response.ERROR, resp.Code)
	assert.NotEmpty(t, resp.Msg)
}

// TestLogOperation_JSONBindingError tests JSON binding error path in LogOperation
func TestLogOperation_JSONBindingError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Send invalid JSON to trigger binding error
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/logOperation", bytes.NewReader([]byte("invalid-json")))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	LogOperation(c)

	// Assert - HTTP status is always 200
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Should return error code due to JSON binding failure
	assert.Equal(t, response.ERROR, resp.Code)
	assert.NotEmpty(t, resp.Msg)
}

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
		checkMsg     func(msg string) bool
	}{
		{
			name:         "empty token",
			requestBody:  map[string]string{"token": ""},
			expectedCode: response.ERROR,
			checkMsg: func(msg string) bool {
				return msg != ""
			},
		},
		{
			name: "invalid token format",
			requestBody: map[string]string{
				"token": "invalid-token-string",
			},
			expectedCode: response.ERROR,
			checkMsg: func(msg string) bool {
				return msg != ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/validateToken", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute
			ValidateToken(c)

			// Assert - HTTP status is always 200
			assert.Equal(t, http.StatusOK, w.Code)

			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedCode, resp.Code)
			if tt.checkMsg != nil {
				assert.True(t, tt.checkMsg(resp.Msg))
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
		checkMsg     func(msg string) bool
	}{
		{
			name: "empty token",
			requestBody: map[string]interface{}{
				"token":     "",
				"expiresAt": 1700000000,
			},
			expectedCode: response.ERROR,
			checkMsg: func(msg string) bool {
				return msg != ""
			},
		},
		{
			name: "invalid token",
			requestBody: map[string]interface{}{
				"token":     "invalid-token",
				"expiresAt": 1700000000,
			},
			expectedCode: response.ERROR,
			checkMsg: func(msg string) bool {
				return msg != ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/refreshToken", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute
			RefreshToken(c)

			// Assert - HTTP status is always 200
			assert.Equal(t, http.StatusOK, w.Code)

			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedCode, resp.Code)
			if tt.checkMsg != nil {
				assert.True(t, tt.checkMsg(resp.Msg))
			}
		})
	}
}

func TestCheckPermission(t *testing.T) {
	// Note: This test requires database and Casbin initialization.
	// To enable this test, you need to:
	// 1. Initialize global.GVA_CONFIG with valid MySQL configuration
	// 2. Ensure Casbin model file exists at global.GVA_CONFIG.Casbin.ModelPath
	// 3. Initialize global.GVA_DB with a test database connection
	//
	// Example setup:
	// global.GVA_CONFIG.Mysql = config.Mysql{
	//     Path:     "127.0.0.1",
	//     Port:     "3306",
	//     Config:   "charset=utf8mb4&parseTime=True&loc=Local",
	//     Dbname:   "test_db",
	//     Username: "root",
	//     Password: "password",
	// }
	// global.GVA_CONFIG.Casbin = config.Casbin{
	//     ModelPath: "./resource/rbac_model.conf",
	// }
	t.Skip("Skipping: requires database and casbin initialization")

	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
		expectedData map[string]interface{}
	}{
		{
			name: "valid permission check - allowed",
			requestBody: map[string]string{
				"sub": "admin",
				"obj": "/api/users",
				"act": "GET",
			},
			expectedCode: response.SUCCESS,
			expectedData: map[string]interface{}{"allowed": true},
		},
		{
			name: "valid permission check - denied",
			requestBody: map[string]string{
				"sub": "guest",
				"obj": "/api/admin",
				"act": "DELETE",
			},
			expectedCode: response.SUCCESS,
			expectedData: map[string]interface{}{"allowed": false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/checkPermission", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute
			CheckPermission(c)

			// Assert - HTTP status is always 200
			assert.Equal(t, http.StatusOK, w.Code)

			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedCode, resp.Code)

			// Check response data
			if tt.expectedData != nil {
				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, tt.expectedData["allowed"], data["allowed"])
			}
		})
	}
}

func TestLogOperation(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "valid operation log",
			requestBody: map[string]interface{}{
				"ip":      "127.0.0.1",
				"method":  "POST",
				"path":    "/api/users",
				"status":  200,
				"latency": 1000000,
				"agent":   "test-agent",
				"user_id": 1,
			},
			expectedCode: response.SUCCESS,
			expectedMsg:  "operation logged successfully",
		},
		{
			name: "operation log with error message",
			requestBody: map[string]interface{}{
				"ip":            "192.168.1.1",
				"method":        "DELETE",
				"path":          "/api/users/1",
				"status":        500,
				"latency":       5000000,
				"agent":         "curl/7.68.0",
				"error_message": "internal server error",
				"user_id":       2,
			},
			expectedCode: response.SUCCESS,
			expectedMsg:  "operation logged successfully",
		},
		{
			name: "minimal operation log",
			requestBody: map[string]interface{}{
				"ip":     "10.0.0.1",
				"method": "GET",
				"path":   "/health",
			},
			expectedCode: response.SUCCESS,
			expectedMsg:  "operation logged successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip if database is not initialized (unit test environment)
			if global.GVA_DB == nil {
				t.Skip("Skipping test: GVA_DB not initialized in unit test environment")
			}

			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/logOperation", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute
			LogOperation(c)

			// Assert - HTTP status is always 200
			assert.Equal(t, http.StatusOK, w.Code)

			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedCode, resp.Code)
			if tt.expectedMsg != "" {
				assert.Equal(t, tt.expectedMsg, resp.Msg)
			}
		})
	}
}
