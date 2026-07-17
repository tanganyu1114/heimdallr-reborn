package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tanganyu1114/heimdallr-reborn/server/global"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	"github.com/tanganyu1114/heimdallr-reborn/server/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func init() {
	// Initialize global logger for tests
	global.GVA_LOG = zap.NewNop()
}

func TestGetPublicKey(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "missing captchaId",
			requestBody:  map[string]string{},
			expectedCode: response.ERROR,
			expectedMsg:  "请求参数错误",
		},
		{
			name: "valid captchaId",
			requestBody: map[string]string{
				"captchaId": "test-captcha-id",
			},
			expectedCode: response.SUCCESS,
			expectedMsg:  "获取成功",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/base/publicKey", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute
			GetPublicKey(c)

			// Assert - HTTP status is always 200
			assert.Equal(t, http.StatusOK, w.Code)

			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedCode, resp.Code)
			if tt.expectedMsg != "" {
				assert.Equal(t, tt.expectedMsg, resp.Msg)
			}

			if tt.expectedCode == response.SUCCESS {
				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.NotEmpty(t, data["publicKey"])
				assert.NotEmpty(t, data["challenge"])
			}
		})
	}
}

func TestGetSDKChallenge(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "missing apiKey",
			requestBody:  map[string]string{},
			expectedCode: response.ERROR,
			expectedMsg:  "请求参数错误",
		},
		{
			name: "valid apiKey",
			requestBody: map[string]string{
				"apiKey": "test-api-key",
			},
			expectedCode: response.SUCCESS,
			expectedMsg:  "获取成功",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/base/sdkChallenge", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute
			GetSDKChallenge(c)

			// Assert - HTTP status is always 200
			assert.Equal(t, http.StatusOK, w.Code)

			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedCode, resp.Code)
			if tt.expectedMsg != "" {
				assert.Equal(t, tt.expectedMsg, resp.Msg)
			}

			if tt.expectedCode == response.SUCCESS {
				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.NotEmpty(t, data["publicKey"])
				assert.NotEmpty(t, data["challenge"])
			}
		})
	}
}

func TestLogin(t *testing.T) {
	// Generate RSA keys first
	utils.GenerateRSAKeys()

	// Get public key for encryption
	publicKey, _, _ := utils.GetPublicKeyWithChallenge("test-captcha-id")

	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "missing encrypted_data",
			requestBody:  map[string]string{},
			expectedCode: response.ERROR,
			expectedMsg:  "请求参数错误，请检查加密数据",
		},
		{
			name: "invalid encrypted data",
			requestBody: map[string]string{
				"encrypted_data": "invalid-base64-data",
			},
			expectedCode: response.ERROR,
			expectedMsg:  "登录数据解密失败，请重试",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/base/login", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute
			Login(c)

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

	// Test valid encrypted login data separately
	// Note: This test will fail because captcha and database are not mocked
	// It's kept for reference but skipped in CI
	t.Run("valid encrypted login data", func(t *testing.T) {
		t.Skip("Skipping: requires database and captcha mock")

		// Create a valid login request
		loginData := request.Login{
			Username:  "testuser",
			Password:  "testpass",
			Captcha:   "123456",
			CaptchaId: "test-captcha-id",
		}

		// Add challenge
		_, challenge, _ := utils.GetPublicKeyWithChallenge(loginData.CaptchaId)
		loginData.Challenge = challenge

		// Encrypt login data
		loginJSON, _ := json.Marshal(loginData)
		encryptedData, _ := utils.RSAEncrypt(publicKey, string(loginJSON))

		// Setup
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody := map[string]string{
			"encrypted_data": encryptedData,
		}
		body, _ := json.Marshal(requestBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/base/login", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Execute
		Login(c)

		// Assert - will fail because captcha verification fails, but that's expected
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestSDKLogin(t *testing.T) {
	// Generate RSA keys first
	utils.GenerateRSAKeys()

	// Get public key for encryption
	publicKey, _, _ := utils.GetPublicKeyWithChallenge("test-api-key")

	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "missing encrypted_data",
			requestBody:  map[string]string{},
			expectedCode: response.ERROR,
			expectedMsg:  "请求参数错误，请检查加密数据",
		},
		{
			name: "invalid encrypted data",
			requestBody: map[string]string{
				"encrypted_data": "invalid-base64-data",
			},
			expectedCode: response.ERROR,
			expectedMsg:  "SDK登录数据解密失败，请重试",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/base/sdkLogin", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute
			SDKLogin(c)

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

	// Test valid encrypted SDK login data separately
	// Note: This test will fail because database is not mocked
	// It's kept for reference but skipped in CI
	t.Run("valid encrypted SDK login data", func(t *testing.T) {
		t.Skip("Skipping: requires database mock")

		// Create a valid SDK login request
		sdkLoginData := request.SDKLogin{
			APIKey:    "test-api-key",
			APISecret: "test-api-secret",
		}

		// Add challenge
		_, challenge, _ := utils.GetPublicKeyWithChallenge(sdkLoginData.APIKey)
		sdkLoginData.Challenge = challenge

		// Encrypt SDK login data
		sdkLoginJSON, _ := json.Marshal(sdkLoginData)
		encryptedData, _ := utils.RSAEncrypt(publicKey, string(sdkLoginJSON))

		// Setup
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody := map[string]string{
			"encrypted_data": encryptedData,
		}
		body, _ := json.Marshal(requestBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/base/sdkLogin", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Execute
		SDKLogin(c)

		// Assert - will fail because API key verification fails, but that's expected
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
