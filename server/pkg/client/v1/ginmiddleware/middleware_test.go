package ginmiddleware

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	clientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/service"
	"go.uber.org/mock/gomock"
)

// args represents the arguments for table-driven tests
type args struct {
	c        *gin.Context
	err      error
	data     gin.H
	msg      string
	b        []byte
	client   *clientv1.Client
	apiGroup string
}

func TestGinMiddleware_CasbinHandler(t *testing.T) {
	t.Run("returns gin.HandlerFunc", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}
		handler := gm.CasbinHandler()
		assert.NotNil(t, handler)
	})

	t.Run("handler denies access when claims missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		handler := gm.CasbinHandler()
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("handler allows access when permission granted", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockAuthzService.EXPECT().CheckPermission("888", "/test", http.MethodGet).Return(&response.CheckPermissionResponse{
			Allowed: true,
		}, nil)

		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		c.Set(ClaimsKey, &request.CustomClaims{
			ID:          1,
			AuthorityId: "888",
		})

		handler := gm.CasbinHandler()
		handler(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("handler denies access when permission check fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockAuthzService.EXPECT().CheckPermission("888", "/test", http.MethodGet).Return(nil, fmt.Errorf("permission check error"))

		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		c.Set(ClaimsKey, &request.CustomClaims{
			ID:          1,
			AuthorityId: "888",
		})

		handler := gm.CasbinHandler()
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("handler denies access when permission not granted", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockAuthzService.EXPECT().CheckPermission("888", "/test", http.MethodGet).Return(&response.CheckPermissionResponse{
			Allowed: false,
		}, nil)

		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		c.Set(ClaimsKey, &request.CustomClaims{
			ID:          1,
			AuthorityId: "888",
		})

		handler := gm.CasbinHandler()
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestGinMiddleware_JWTAuth(t *testing.T) {
	t.Run("returns gin.HandlerFunc", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}
		handler := gm.JWTAuth()
		assert.NotNil(t, handler)
	})

	t.Run("handler denies access when token missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		handler := gm.JWTAuth()
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("handler allows access with valid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockAuthzService.EXPECT().ValidateToken("valid-token").Return(&response.ValidateTokenResponse{
			Claims: &request.CustomClaims{
				ID:          1,
				AuthorityId: "888",
				BufferTime:  86400,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Unix() + 100000,
				},
			},
		}, nil)

		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(TokenHeader, "valid-token")
		c.Request = req

		handler := gm.JWTAuth()
		handler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		claims, exists := c.Get(ClaimsKey)
		assert.True(t, exists)
		assert.NotNil(t, claims)
	})

	t.Run("handler denies access when token validation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockAuthzService.EXPECT().ValidateToken("invalid-token").Return(nil, fmt.Errorf("invalid token"))

		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(TokenHeader, "invalid-token")
		c.Request = req

		handler := gm.JWTAuth()
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("handler refreshes token when close to expiry", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockAuthzService.EXPECT().ValidateToken("expiring-token").Return(&response.ValidateTokenResponse{
			Claims: &request.CustomClaims{
				ID:          1,
				AuthorityId: "888",
				BufferTime:  86400,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Unix() + 3600,
				},
			},
		}, nil)
		mockAuthzService.EXPECT().RefreshToken(gomock.Any(), gomock.Any()).Return(&response.RefreshTokenResponse{
			Token:     "new-token",
			ExpiresAt: time.Now().Unix() + 60*60*24*7,
		}, nil)

		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(TokenHeader, "expiring-token")
		c.Request = req

		handler := gm.JWTAuth()
		handler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Header().Get(NewTokenHeader))
		assert.NotEmpty(t, w.Header().Get(NewExpiresAtHeader))
	})
}

func TestGinMiddleware_OperationRecord(t *testing.T) {
	t.Run("returns gin.HandlerFunc", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockAuthzService.EXPECT().LogOperation(gomock.Any()).Return(nil).AnyTimes()

		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}
		handler := gm.OperationRecord()
		assert.NotNil(t, handler)
	})

	t.Run("handler logs operation for POST request", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockAuthzService.EXPECT().LogOperation(gomock.Any()).Return(nil)

		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := bytes.NewBufferString(`{"key":"value"}`)
		req := httptest.NewRequest(http.MethodPost, "/test", body)
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Set(ClaimsKey, &request.CustomClaims{ID: 1})

		handler := gm.OperationRecord()
		handler(c)

		assert.Equal(t, http.StatusOK, w.Code)

		// Wait for goroutine to complete
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("handler logs SDK login operation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockAuthzService.EXPECT().LogOperation(gomock.Any()).Do(func(record *request.LogOperationRequest) {
			assert.True(t, record.IsSDKLogin)
		}).Return(nil)

		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Set(ClaimsKey, &request.CustomClaims{ID: 1})
		c.Set(SDKLoginKey, true)

		handler := gm.OperationRecord()
		handler(c)

		assert.Equal(t, http.StatusOK, w.Code)

		// Wait for goroutine to complete
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("handler logs binary response without body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockAuthzService.EXPECT().LogOperation(gomock.Any()).Do(func(record *request.LogOperationRequest) {
			assert.Equal(t, "[Binary File - Not Logged]", record.Resp)
		}).Return(nil)

		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		c.Request = req
		c.Set(ClaimsKey, &request.CustomClaims{ID: 1})

		// Register route with binary response handler
		r.Use(gm.OperationRecord())
		r.GET("/test", func(ctx *gin.Context) {
			ctx.Header("Content-Type", "application/pdf")
			ctx.Writer.Write([]byte("binary content"))
		})

		// Execute through router to properly trigger middleware chain
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Wait for goroutine to complete
		time.Sleep(10 * time.Millisecond)
	})
}

func TestGinMiddleware_logOperation(t *testing.T) {
	t.Run("logs operation successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockAuthzService.EXPECT().LogOperation(gomock.Any()).Return(nil)

		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}

		record := &request.LogOperationRequest{
			IP:     "127.0.0.1",
			Method: http.MethodPost,
			Path:   "/test",
			Body:   `{"key":"value"}`,
			UserID: 1,
		}

		assert.NotPanics(t, func() {
			gm.logOperation(record)
		})
	})

	t.Run("logs warning when log operation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockAuthzService.EXPECT().LogOperation(gomock.Any()).Return(fmt.Errorf("log service error"))

		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}

		record := &request.LogOperationRequest{
			IP:     "127.0.0.1",
			Method: http.MethodPost,
			Path:   "/test",
			Body:   `{"key":"value"}`,
			UserID: 1,
		}

		// Should not panic even when logging fails
		assert.NotPanics(t, func() {
			gm.logOperation(record)
		})
	})
}

func TestGinMiddleware_refreshToken(t *testing.T) {
	t.Run("refreshes token successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockAuthzService.EXPECT().RefreshToken(gomock.Any(), gomock.Any()).Return(&response.RefreshTokenResponse{
			Token:     "new-token",
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		}, nil)

		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		claims := &request.CustomClaims{
			ID:          1,
			AuthorityId: "888",
			BufferTime:  86400,
		}

		assert.NotPanics(t, func() {
			gm.refreshToken(c, "old-token", claims)
		})

		assert.NotEmpty(t, w.Header().Get(NewTokenHeader))
	})

	t.Run("logs error when refresh token fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockAuthzService.EXPECT().RefreshToken(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("refresh service error"))

		mockFactory := service.NewMockFactory(ctrl)
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		gm := &GinMiddleware{
			client:   &clientv1.Client{Factory: mockFactory},
			apiGroup: "test",
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		claims := &request.CustomClaims{
			ID:          1,
			AuthorityId: "888",
			BufferTime:  86400,
		}

		// Should not panic even when refresh fails
		assert.NotPanics(t, func() {
			gm.refreshToken(c, "old-token", claims)
		})

		// No new token header should be set
		assert.Empty(t, w.Header().Get(NewTokenHeader))
	})
}

func TestGinMiddleware_shouldRefreshToken(t *testing.T) {
	tests := []struct {
		name   string
		claims *request.CustomClaims
		want   bool
	}{
		{
			name: "should refresh when close to expiry",
			claims: &request.CustomClaims{
				BufferTime: 86400,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Unix() + 3600,
				},
			},
			want: true,
		},
		{
			name: "should not refresh when far from expiry",
			claims: &request.CustomClaims{
				BufferTime: 86400,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Unix() + 100000,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gm := &GinMiddleware{
				apiGroup: "test",
			}
			assert.Equal(t, tt.want, gm.shouldRefreshToken(tt.claims))
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		args    args
		want    *GinMiddleware
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "creates GinMiddleware successfully",
			args: args{
				client:   &clientv1.Client{},
				apiGroup: "test",
			},
			want: &GinMiddleware{
				apiGroup: "test",
			},
			wantErr: assert.NoError,
		},
		{
			name: "returns error when client is nil",
			args: args{
				client:   nil,
				apiGroup: "test",
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "returns error when apiGroup is empty",
			args: args{
				client:   &clientv1.Client{},
				apiGroup: "",
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.client, tt.args.apiGroup)
			if !tt.wantErr(t, err, fmt.Sprintf("New(%v, %v)", tt.args.client, tt.args.apiGroup)) {
				return
			}
			if tt.want != nil {
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.apiGroup, got.apiGroup)
			}
		})
	}
}

func Test_extractError(t *testing.T) {
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "returns empty string when no errors",
			args: args{
				c: func() *gin.Context {
					gin.SetMode(gin.TestMode)
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
					return c
				}(),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, extractError(tt.args.c))
		})
	}
}

func Test_extractUserID(t *testing.T) {
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "extracts user ID from claims",
			args: args{
				c: func() *gin.Context {
					gin.SetMode(gin.TestMode)
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
					c.Set(ClaimsKey, &request.CustomClaims{ID: 123})
					return c
				}(),
			},
			want: 123,
		},
		{
			name: "extracts user ID from header when claims missing",
			args: args{
				c: func() *gin.Context {
					gin.SetMode(gin.TestMode)
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					req := httptest.NewRequest(http.MethodGet, "/test", nil)
					req.Header.Set(UserIDHeader, "456")
					c.Request = req
					return c
				}(),
			},
			want: 456,
		},
		{
			name: "returns 0 when no claims and no header",
			args: args{
				c: func() *gin.Context {
					gin.SetMode(gin.TestMode)
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
					return c
				}(),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, extractUserID(tt.args.c))
		})
	}
}

func Test_handleTokenError(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{
			name: "handles expired token error",
			args: args{
				c: func() *gin.Context {
					gin.SetMode(gin.TestMode)
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
					return c
				}(),
				err: fmt.Errorf("token expired"),
			},
		},
		{
			name: "handles generic token error",
			args: args{
				c: func() *gin.Context {
					gin.SetMode(gin.TestMode)
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
					return c
				}(),
				err: fmt.Errorf("invalid token"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleTokenError(tt.args.c, tt.args.err)
		})
	}
}

func Test_isBinaryResponse(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		want        bool
	}{
		{
			name:        "detects Excel content type",
			contentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			want:        true,
		},
		{
			name:        "detects octet stream",
			contentType: "application/octet-stream",
			want:        true,
		},
		{
			name:        "detects PDF content type",
			contentType: "application/pdf",
			want:        true,
		},
		{
			name:        "detects image content type",
			contentType: "image/png",
			want:        true,
		},
		{
			name:        "returns false for JSON",
			contentType: "application/json",
			want:        false,
		},
		{
			name:        "returns false for HTML",
			contentType: "text/html",
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isBinaryResponse(tt.contentType))
		})
	}
}

func Test_readBody(t *testing.T) {
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "reads request body",
			args: args{
				c: func() *gin.Context {
					gin.SetMode(gin.TestMode)
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					body := bytes.NewBufferString(`{"key":"value"}`)
					req := httptest.NewRequest(http.MethodPost, "/test", body)
					c.Request = req
					return c
				}(),
			},
			want: []byte(`{"key":"value"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := readBody(tt.args.c)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_respondError(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{
			name: "responds with error JSON",
			args: args{
				c: func() *gin.Context {
					gin.SetMode(gin.TestMode)
					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
					return c
				}(),
				data: gin.H{"reload": true},
				msg:  "未登录或非法访问",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respondError(tt.args.c, tt.args.data, tt.args.msg)
		})
	}
}
