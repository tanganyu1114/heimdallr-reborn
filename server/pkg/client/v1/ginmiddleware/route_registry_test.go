package ginmiddleware

import (
	"net/http"
	"reflect"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tanganyu1114/heimdallr-reborn/server/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	clientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/service"
	"go.uber.org/mock/gomock"
)

// MockFactory is an alias for service.MockFactory for convenience
type MockFactory = service.MockFactory

// NewMockFactory creates a new MockFactory instance
func NewMockFactory(ctrl *gomock.Controller) *MockFactory {
	return service.NewMockFactory(ctrl)
}

func TestGinMiddleware_Group(t *testing.T) {
	type fields struct {
		client   *clientv1.Client
		apiGroup string
	}
	type args struct {
		rg           *gin.RouterGroup
		relativePath string
		handlers     []gin.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRegistry
	}{
		{
			name: "create route registry with valid group",
			fields: fields{
				client:   nil,
				apiGroup: "test-api-group",
			},
			args: args{
				rg:           &gin.RouterGroup{},
				relativePath: "/api",
				handlers:     nil,
			},
			want: &RouteRegistry{
				apiGroup: "test-api-group",
			},
		},
		{
			name: "create route registry with handlers",
			fields: fields{
				client:   nil,
				apiGroup: "admin",
			},
			args: args{
				rg:           &gin.RouterGroup{},
				relativePath: "/admin",
				handlers:     []gin.HandlerFunc{func(c *gin.Context) {}},
			},
			want: &RouteRegistry{
				apiGroup: "admin",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gm := &GinMiddleware{
				client:   tt.fields.client,
				apiGroup: tt.fields.apiGroup,
			}
			got := gm.Group(tt.args.rg, tt.args.relativePath, tt.args.handlers...)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.apiGroup, got.apiGroup)
			assert.Equal(t, tt.fields.client, got.client)
		})
	}
}

func TestGinMiddleware_Group_PanicsOnNilRouterGroup(t *testing.T) {
	gm := &GinMiddleware{
		client:   nil,
		apiGroup: "test",
	}
	assert.Panics(t, func() {
		gm.Group(nil, "/api")
	})
}

func TestRouteRegistry_Any(t *testing.T) {
	type fields struct {
		group      *gin.RouterGroup
		client     *clientv1.Client
		apiGroup   string
		mu         sync.Mutex
		apiCache   map[string]*model.SysApi
		cacheOnce  sync.Once
		strictMode bool
	}
	type args struct {
		relativePath string
		desc         string
		handlers     []gin.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRegistry
	}{
		{
			name: "register ANY route with nil client",
			fields: fields{
				client:     nil,
				apiGroup:   "test-group",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				relativePath: "/test",
				desc:         "Test endpoint",
				handlers:     []gin.HandlerFunc{func(c *gin.Context) {}},
			},
			want: &RouteRegistry{
				apiGroup: "test-group",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			engine := gin.New()
			tt.fields.group = engine.Group("/api")

			a := &RouteRegistry{
				group:      tt.fields.group,
				client:     tt.fields.client,
				apiGroup:   tt.fields.apiGroup,
				mu:         tt.fields.mu,
				apiCache:   tt.fields.apiCache,
				cacheOnce:  tt.fields.cacheOnce,
				strictMode: tt.fields.strictMode,
			}
			got := a.Any(tt.args.relativePath, tt.args.desc, tt.args.handlers...)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.apiGroup, got.apiGroup)
			assert.Equal(t, tt.fields.client, got.client)
		})
	}
}

func TestRouteRegistry_DELETE(t *testing.T) {
	type fields struct {
		group      *gin.RouterGroup
		client     *clientv1.Client
		apiGroup   string
		mu         sync.Mutex
		apiCache   map[string]*model.SysApi
		cacheOnce  sync.Once
		strictMode bool
	}
	type args struct {
		relativePath string
		desc         string
		handlers     []gin.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRegistry
	}{
		{
			name: "register DELETE route with nil client",
			fields: fields{
				client:     nil,
				apiGroup:   "user-api",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				relativePath: "/users/:id",
				desc:         "Delete user",
				handlers:     []gin.HandlerFunc{func(c *gin.Context) {}},
			},
			want: &RouteRegistry{
				apiGroup: "user-api",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			engine := gin.New()
			tt.fields.group = engine.Group("/api")

			a := &RouteRegistry{
				group:      tt.fields.group,
				client:     tt.fields.client,
				apiGroup:   tt.fields.apiGroup,
				mu:         tt.fields.mu,
				apiCache:   tt.fields.apiCache,
				cacheOnce:  tt.fields.cacheOnce,
				strictMode: tt.fields.strictMode,
			}
			got := a.DELETE(tt.args.relativePath, tt.args.desc, tt.args.handlers...)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.apiGroup, got.apiGroup)
			assert.Equal(t, tt.fields.client, got.client)
		})
	}
}

func TestRouteRegistry_GET(t *testing.T) {
	type fields struct {
		group      *gin.RouterGroup
		client     *clientv1.Client
		apiGroup   string
		mu         sync.Mutex
		apiCache   map[string]*model.SysApi
		cacheOnce  sync.Once
		strictMode bool
	}
	type args struct {
		relativePath string
		desc         string
		handlers     []gin.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRegistry
	}{
		{
			name: "register GET route with nil client",
			fields: fields{
				client:     nil,
				apiGroup:   "user-api",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				relativePath: "/users",
				desc:         "Get user list",
				handlers:     []gin.HandlerFunc{func(c *gin.Context) {}},
			},
			want: &RouteRegistry{
				apiGroup: "user-api",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			engine := gin.New()
			tt.fields.group = engine.Group("/api")

			a := &RouteRegistry{
				group:      tt.fields.group,
				client:     tt.fields.client,
				apiGroup:   tt.fields.apiGroup,
				mu:         tt.fields.mu,
				apiCache:   tt.fields.apiCache,
				cacheOnce:  tt.fields.cacheOnce,
				strictMode: tt.fields.strictMode,
			}
			got := a.GET(tt.args.relativePath, tt.args.desc, tt.args.handlers...)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.apiGroup, got.apiGroup)
			assert.Equal(t, tt.fields.client, got.client)
		})
	}
}

func TestRouteRegistry_Group(t *testing.T) {
	type fields struct {
		client     *clientv1.Client
		apiGroup   string
		mu         sync.Mutex
		apiCache   map[string]*model.SysApi
		cacheOnce  sync.Once
		strictMode bool
	}
	type args struct {
		relativePath string
		handlers     []gin.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRegistry
	}{
		{
			name: "create sub-group without handlers",
			fields: fields{
				client:     nil,
				apiGroup:   "v1",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				relativePath: "/users",
				handlers:     nil,
			},
			want: &RouteRegistry{
				apiGroup: "v1",
			},
		},
		{
			name: "create sub-group with middleware",
			fields: fields{
				client:     nil,
				apiGroup:   "v1",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				relativePath: "/admin",
				handlers:     []gin.HandlerFunc{func(c *gin.Context) {}},
			},
			want: &RouteRegistry{
				apiGroup: "v1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			engine := gin.New()
			group := engine.Group("/api")

			a := &RouteRegistry{
				group:      group,
				client:     tt.fields.client,
				apiGroup:   tt.fields.apiGroup,
				mu:         tt.fields.mu,
				apiCache:   tt.fields.apiCache,
				cacheOnce:  tt.fields.cacheOnce,
				strictMode: tt.fields.strictMode,
			}
			got := a.Group(tt.args.relativePath, tt.args.handlers...)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.apiGroup, got.apiGroup)
			assert.Equal(t, tt.fields.client, got.client)
		})
	}
}

func TestRouteRegistry_HEAD(t *testing.T) {
	type fields struct {
		client     *clientv1.Client
		apiGroup   string
		mu         sync.Mutex
		apiCache   map[string]*model.SysApi
		cacheOnce  sync.Once
		strictMode bool
	}
	type args struct {
		relativePath string
		desc         string
		handlers     []gin.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRegistry
	}{
		{
			name: "register HEAD route with nil client",
			fields: fields{
				client:     nil,
				apiGroup:   "health-api",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				relativePath: "/health",
				desc:         "Health check",
				handlers:     []gin.HandlerFunc{func(c *gin.Context) {}},
			},
			want: &RouteRegistry{
				apiGroup: "health-api",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			engine := gin.New()
			group := engine.Group("/api")

			a := &RouteRegistry{
				group:      group,
				client:     tt.fields.client,
				apiGroup:   tt.fields.apiGroup,
				mu:         tt.fields.mu,
				apiCache:   tt.fields.apiCache,
				cacheOnce:  tt.fields.cacheOnce,
				strictMode: tt.fields.strictMode,
			}
			got := a.HEAD(tt.args.relativePath, tt.args.desc, tt.args.handlers...)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.apiGroup, got.apiGroup)
			assert.Equal(t, tt.fields.client, got.client)
		})
	}
}

func TestRouteRegistry_OPTIONS(t *testing.T) {
	type fields struct {
		client     *clientv1.Client
		apiGroup   string
		mu         sync.Mutex
		apiCache   map[string]*model.SysApi
		cacheOnce  sync.Once
		strictMode bool
	}
	type args struct {
		relativePath string
		desc         string
		handlers     []gin.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRegistry
	}{
		{
			name: "register OPTIONS route with nil client",
			fields: fields{
				client:     nil,
				apiGroup:   "cors-api",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				relativePath: "/cors",
				desc:         "CORS preflight",
				handlers:     []gin.HandlerFunc{func(c *gin.Context) {}},
			},
			want: &RouteRegistry{
				apiGroup: "cors-api",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			engine := gin.New()
			group := engine.Group("/api")

			a := &RouteRegistry{
				group:      group,
				client:     tt.fields.client,
				apiGroup:   tt.fields.apiGroup,
				mu:         tt.fields.mu,
				apiCache:   tt.fields.apiCache,
				cacheOnce:  tt.fields.cacheOnce,
				strictMode: tt.fields.strictMode,
			}
			got := a.OPTIONS(tt.args.relativePath, tt.args.desc, tt.args.handlers...)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.apiGroup, got.apiGroup)
			assert.Equal(t, tt.fields.client, got.client)
		})
	}
}

func TestRouteRegistry_PATCH(t *testing.T) {
	type fields struct {
		client     *clientv1.Client
		apiGroup   string
		mu         sync.Mutex
		apiCache   map[string]*model.SysApi
		cacheOnce  sync.Once
		strictMode bool
	}
	type args struct {
		relativePath string
		desc         string
		handlers     []gin.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRegistry
	}{
		{
			name: "register PATCH route with nil client",
			fields: fields{
				client:     nil,
				apiGroup:   "user-api",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				relativePath: "/users/:id",
				desc:         "Update user partially",
				handlers:     []gin.HandlerFunc{func(c *gin.Context) {}},
			},
			want: &RouteRegistry{
				apiGroup: "user-api",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			engine := gin.New()
			group := engine.Group("/api")

			a := &RouteRegistry{
				group:      group,
				client:     tt.fields.client,
				apiGroup:   tt.fields.apiGroup,
				mu:         tt.fields.mu,
				apiCache:   tt.fields.apiCache,
				cacheOnce:  tt.fields.cacheOnce,
				strictMode: tt.fields.strictMode,
			}
			got := a.PATCH(tt.args.relativePath, tt.args.desc, tt.args.handlers...)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.apiGroup, got.apiGroup)
			assert.Equal(t, tt.fields.client, got.client)
		})
	}
}

func TestRouteRegistry_POST(t *testing.T) {
	type fields struct {
		client     *clientv1.Client
		apiGroup   string
		mu         sync.Mutex
		apiCache   map[string]*model.SysApi
		cacheOnce  sync.Once
		strictMode bool
	}
	type args struct {
		relativePath string
		desc         string
		handlers     []gin.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRegistry
	}{
		{
			name: "register POST route with nil client",
			fields: fields{
				client:     nil,
				apiGroup:   "user-api",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				relativePath: "/users",
				desc:         "Create user",
				handlers:     []gin.HandlerFunc{func(c *gin.Context) {}},
			},
			want: &RouteRegistry{
				apiGroup: "user-api",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			engine := gin.New()
			group := engine.Group("/api")

			a := &RouteRegistry{
				group:      group,
				client:     tt.fields.client,
				apiGroup:   tt.fields.apiGroup,
				mu:         tt.fields.mu,
				apiCache:   tt.fields.apiCache,
				cacheOnce:  tt.fields.cacheOnce,
				strictMode: tt.fields.strictMode,
			}
			got := a.POST(tt.args.relativePath, tt.args.desc, tt.args.handlers...)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.apiGroup, got.apiGroup)
			assert.Equal(t, tt.fields.client, got.client)
		})
	}
}

func TestRouteRegistry_PUT(t *testing.T) {
	type fields struct {
		client     *clientv1.Client
		apiGroup   string
		mu         sync.Mutex
		apiCache   map[string]*model.SysApi
		cacheOnce  sync.Once
		strictMode bool
	}
	type args struct {
		relativePath string
		desc         string
		handlers     []gin.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRegistry
	}{
		{
			name: "register PUT route with nil client",
			fields: fields{
				client:     nil,
				apiGroup:   "user-api",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				relativePath: "/users/:id",
				desc:         "Update user",
				handlers:     []gin.HandlerFunc{func(c *gin.Context) {}},
			},
			want: &RouteRegistry{
				apiGroup: "user-api",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			engine := gin.New()
			group := engine.Group("/api")

			a := &RouteRegistry{
				group:      group,
				client:     tt.fields.client,
				apiGroup:   tt.fields.apiGroup,
				mu:         tt.fields.mu,
				apiCache:   tt.fields.apiCache,
				cacheOnce:  tt.fields.cacheOnce,
				strictMode: tt.fields.strictMode,
			}
			got := a.PUT(tt.args.relativePath, tt.args.desc, tt.args.handlers...)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.apiGroup, got.apiGroup)
			assert.Equal(t, tt.fields.client, got.client)
		})
	}
}

func TestRouteRegistry_Use(t *testing.T) {
	type fields struct {
		client     *clientv1.Client
		apiGroup   string
		mu         sync.Mutex
		apiCache   map[string]*model.SysApi
		cacheOnce  sync.Once
		strictMode bool
	}
	type args struct {
		middleware []gin.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRegistry
	}{
		{
			name: "add single middleware",
			fields: fields{
				client:     nil,
				apiGroup:   "test-api",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				middleware: []gin.HandlerFunc{func(c *gin.Context) {}},
			},
			want: &RouteRegistry{
				apiGroup: "test-api",
			},
		},
		{
			name: "add multiple middleware",
			fields: fields{
				client:     nil,
				apiGroup:   "test-api",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				middleware: []gin.HandlerFunc{
					func(c *gin.Context) {},
					func(c *gin.Context) {},
				},
			},
			want: &RouteRegistry{
				apiGroup: "test-api",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			engine := gin.New()
			group := engine.Group("/api")

			a := &RouteRegistry{
				group:      group,
				client:     tt.fields.client,
				apiGroup:   tt.fields.apiGroup,
				mu:         tt.fields.mu,
				apiCache:   tt.fields.apiCache,
				cacheOnce:  tt.fields.cacheOnce,
				strictMode: tt.fields.strictMode,
			}
			got := a.Use(tt.args.middleware...)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.apiGroup, got.apiGroup)
			assert.True(t, reflect.DeepEqual(a, got), "Use should return same instance")
		})
	}
}

func TestRouteRegistry_applyAPI(t *testing.T) {
	type fields struct {
		client     *clientv1.Client
		apiGroup   string
		mu         sync.Mutex
		apiCache   map[string]*model.SysApi
		cacheOnce  sync.Once
		strictMode bool
	}
	type args struct {
		api *model.SysApi
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "nil API object",
			fields: fields{
				client:     nil,
				apiGroup:   "test",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				api: nil,
			},
			wantErr: true,
		},
		{
			name: "empty path",
			fields: fields{
				client:     nil,
				apiGroup:   "test",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				api: &model.SysApi{
					Path:        "",
					Description: "Test API",
					ApiGroup:    "test",
					Method:      "GET",
				},
			},
			wantErr: true,
		},
		{
			name: "empty API group",
			fields: fields{
				client:     nil,
				apiGroup:   "test",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				api: &model.SysApi{
					Path:        "/test",
					Description: "Test API",
					ApiGroup:    "",
					Method:      "GET",
				},
			},
			wantErr: true,
		},
		{
			name: "empty method",
			fields: fields{
				client:     nil,
				apiGroup:   "test",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				api: &model.SysApi{
					Path:        "/test",
					Description: "Test API",
					ApiGroup:    "test",
					Method:      "",
				},
			},
			wantErr: true,
		},
		{
			name: "empty description",
			fields: fields{
				client:     nil,
				apiGroup:   "test",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				api: &model.SysApi{
					Path:        "/test",
					Description: "",
					ApiGroup:    "test",
					Method:      "GET",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			engine := gin.New()
			group := engine.Group("/api")

			a := &RouteRegistry{
				group:      group,
				client:     tt.fields.client,
				apiGroup:   tt.fields.apiGroup,
				mu:         tt.fields.mu,
				apiCache:   tt.fields.apiCache,
				cacheOnce:  tt.fields.cacheOnce,
				strictMode: tt.fields.strictMode,
			}
			err := a.applyAPI(tt.args.api)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRouteRegistry_applyAPIEntry(t *testing.T) {
	type fields struct {
		client     *clientv1.Client
		apiGroup   string
		mu         sync.Mutex
		apiCache   map[string]*model.SysApi
		cacheOnce  sync.Once
		strictMode bool
	}
	type args struct {
		method      string
		subPath     string
		description string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "apply API entry with nil client should not panic",
			fields: fields{
				client:     nil,
				apiGroup:   "test",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				method:      http.MethodGet,
				subPath:     "/test",
				description: "Test endpoint",
			},
		},
		{
			name: "apply API entry with POST method",
			fields: fields{
				client:     nil,
				apiGroup:   "user",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
			args: args{
				method:      http.MethodPost,
				subPath:     "/users",
				description: "Create user",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			engine := gin.New()
			group := engine.Group("/api")

			a := &RouteRegistry{
				group:      group,
				client:     tt.fields.client,
				apiGroup:   tt.fields.apiGroup,
				mu:         tt.fields.mu,
				apiCache:   tt.fields.apiCache,
				cacheOnce:  tt.fields.cacheOnce,
				strictMode: tt.fields.strictMode,
			}
			assert.NotPanics(t, func() {
				a.applyAPIEntry(tt.args.method, tt.args.subPath, tt.args.description)
			})
		})
	}
}

func TestRouteRegistry_applyAPIEntry_WithValidClient(t *testing.T) {
	tests := []struct {
		name       string
		apiGroup   string
		method     string
		subPath    string
		desc       string
		strictMode bool
	}{
		{
			name:       "apply API entry with valid client and GET method",
			apiGroup:   "user",
			method:     http.MethodGet,
			subPath:    "/users",
			desc:       "Get user list",
			strictMode: false,
		},
		{
			name:       "apply API entry with valid client and POST method",
			apiGroup:   "user",
			method:     http.MethodPost,
			subPath:    "/users",
			desc:       "Create user",
			strictMode: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gin.SetMode(gin.TestMode)
			engine := gin.New()
			group := engine.Group("/api")

			mockApiService := service.NewMockApiService(ctrl)
			mockApiService.EXPECT().GetAllApis().Return(&response.SysAPIListResponse{}, nil)
			mockApiService.EXPECT().CreateApi(gomock.Any()).Return(nil)

			mockFactory := NewMockFactory(ctrl)
			mockFactory.EXPECT().Apis().Return(mockApiService).AnyTimes()

			mockClient := &clientv1.Client{
				Factory: mockFactory,
			}

			a := &RouteRegistry{
				group:      group,
				client:     mockClient,
				apiGroup:   tt.apiGroup,
				mu:         sync.Mutex{},
				apiCache:   make(map[string]*model.SysApi),
				cacheOnce:  sync.Once{},
				strictMode: tt.strictMode,
			}
			assert.NotPanics(t, func() {
				a.applyAPIEntry(tt.method, tt.subPath, tt.desc)
			})
		})
	}
}

func TestRouteRegistry_applyAPIEntry_StrictModePanic(t *testing.T) {
	t.Run("apply API entry with strict mode should panic on error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gin.SetMode(gin.TestMode)
		engine := gin.New()
		group := engine.Group("/api")

		mockApiService := service.NewMockApiService(ctrl)
		mockApiService.EXPECT().GetAllApis().Return(&response.SysAPIListResponse{}, nil)
		mockApiService.EXPECT().CreateApi(gomock.Any()).Return(assert.AnError)

		mockFactory := NewMockFactory(ctrl)
		mockFactory.EXPECT().Apis().Return(mockApiService).AnyTimes()

		mockClient := &clientv1.Client{
			Factory: mockFactory,
		}

		a := &RouteRegistry{
			group:      group,
			client:     mockClient,
			apiGroup:   "user",
			mu:         sync.Mutex{},
			apiCache:   make(map[string]*model.SysApi),
			cacheOnce:  sync.Once{},
			strictMode: true,
		}
		assert.Panics(t, func() {
			a.applyAPIEntry(http.MethodPost, "/users", "Create user")
		})
	})
}

func TestRouteRegistry_applyAPIEntry_NonStrictModeError(t *testing.T) {
	t.Run("apply API entry with non-strict mode should not panic on error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gin.SetMode(gin.TestMode)
		engine := gin.New()
		group := engine.Group("/api")

		mockApiService := service.NewMockApiService(ctrl)
		mockApiService.EXPECT().GetAllApis().Return(&response.SysAPIListResponse{}, nil)
		mockApiService.EXPECT().CreateApi(gomock.Any()).Return(assert.AnError)

		mockFactory := NewMockFactory(ctrl)
		mockFactory.EXPECT().Apis().Return(mockApiService).AnyTimes()

		mockClient := &clientv1.Client{
			Factory: mockFactory,
		}

		a := &RouteRegistry{
			group:      group,
			client:     mockClient,
			apiGroup:   "user",
			mu:         sync.Mutex{},
			apiCache:   make(map[string]*model.SysApi),
			cacheOnce:  sync.Once{},
			strictMode: false,
		}
		assert.NotPanics(t, func() {
			a.applyAPIEntry(http.MethodPost, "/users", "Create user")
		})
	})
}

func TestRouteRegistry_applyAPI_CreateNewAPI(t *testing.T) {
	tests := []struct {
		name    string
		api     *model.SysApi
		wantErr bool
	}{
		{
			name: "create new API successfully",
			api: &model.SysApi{
				Path:        "/api/v1/users",
				Description: "Create user",
				ApiGroup:    "user",
				Method:      "POST",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gin.SetMode(gin.TestMode)
			engine := gin.New()
			group := engine.Group("/api")

			mockApiService := service.NewMockApiService(ctrl)
			mockApiService.EXPECT().GetAllApis().Return(&response.SysAPIListResponse{}, nil)
			mockApiService.EXPECT().CreateApi(tt.api).Return(nil)

			mockFactory := NewMockFactory(ctrl)
			mockFactory.EXPECT().Apis().Return(mockApiService).AnyTimes()

			mockClient := &clientv1.Client{
				Factory: mockFactory,
			}

			a := &RouteRegistry{
				group:      group,
				client:     mockClient,
				apiGroup:   "user",
				mu:         sync.Mutex{},
				apiCache:   make(map[string]*model.SysApi),
				cacheOnce:  sync.Once{},
				strictMode: false,
			}

			err := a.applyAPI(tt.api)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRouteRegistry_applyAPI_UpdateExistingAPI(t *testing.T) {
	tests := []struct {
		name        string
		api         *model.SysApi
		existingAPI *model.SysApi
		wantErr     bool
	}{
		{
			name: "update existing API with different description",
			api: &model.SysApi{
				Path:        "/api/v1/users",
				Description: "Updated user creation",
				ApiGroup:    "user",
				Method:      "POST",
			},
			existingAPI: &model.SysApi{
				Path:        "/api/v1/users",
				Description: "Create user",
				ApiGroup:    "user",
				Method:      "POST",
			},
			wantErr: false,
		},
		{
			name: "skip update when description matches",
			api: &model.SysApi{
				Path:        "/api/v1/users",
				Description: "Create user",
				ApiGroup:    "user",
				Method:      "POST",
			},
			existingAPI: &model.SysApi{
				Path:        "/api/v1/users",
				Description: "Create user",
				ApiGroup:    "user",
				Method:      "POST",
			},
			wantErr: false,
		},
		{
			name: "update API returns error",
			api: &model.SysApi{
				Path:        "/api/v1/users",
				Description: "Updated user creation",
				ApiGroup:    "user",
				Method:      "POST",
			},
			existingAPI: &model.SysApi{
				Path:        "/api/v1/users",
				Description: "Create user",
				ApiGroup:    "user",
				Method:      "POST",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gin.SetMode(gin.TestMode)
			engine := gin.New()
			group := engine.Group("/api")

			mockApiService := service.NewMockApiService(ctrl)
			mockApiService.EXPECT().GetAllApis().Return(&response.SysAPIListResponse{
				Apis: []model.SysApi{*tt.existingAPI},
			}, nil)
			if tt.api.Description != tt.existingAPI.Description {
				if tt.wantErr {
					mockApiService.EXPECT().UpdateApi(tt.api).Return(assert.AnError)
				} else {
					mockApiService.EXPECT().UpdateApi(tt.api).Return(nil)
				}
			}

			mockFactory := NewMockFactory(ctrl)
			mockFactory.EXPECT().Apis().Return(mockApiService).AnyTimes()

			mockClient := &clientv1.Client{
				Factory: mockFactory,
			}

			a := &RouteRegistry{
				group:      group,
				client:     mockClient,
				apiGroup:   "user",
				mu:         sync.Mutex{},
				apiCache:   make(map[string]*model.SysApi),
				cacheOnce:  sync.Once{},
				strictMode: false,
			}

			err := a.applyAPI(tt.api)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRouteRegistry_ensureCacheLoaded_WithValidClient(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "ensure cache loaded successfully",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gin.SetMode(gin.TestMode)
			engine := gin.New()
			group := engine.Group("/api")

			mockApiService := service.NewMockApiService(ctrl)
			mockApiService.EXPECT().GetAllApis().Return(&response.SysAPIListResponse{
				Apis: []model.SysApi{
					{
						Path:        "/api/v1/users",
						Description: "Get users",
						ApiGroup:    "user",
						Method:      "GET",
					},
				},
			}, nil)

			mockFactory := NewMockFactory(ctrl)
			mockFactory.EXPECT().Apis().Return(mockApiService).AnyTimes()

			mockClient := &clientv1.Client{
				Factory: mockFactory,
			}

			a := &RouteRegistry{
				group:      group,
				client:     mockClient,
				apiGroup:   "user",
				mu:         sync.Mutex{},
				apiCache:   make(map[string]*model.SysApi),
				cacheOnce:  sync.Once{},
				strictMode: false,
			}

			assert.NotPanics(t, func() {
				a.ensureCacheLoaded()
			})
			assert.Len(t, a.apiCache, 1)
		})
	}
}

func TestRouteRegistry_ensureCacheLoaded_WithNilResponse(t *testing.T) {
	t.Run("ensure cache loaded with nil response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gin.SetMode(gin.TestMode)
		engine := gin.New()
		group := engine.Group("/api")

		mockApiService := service.NewMockApiService(ctrl)
		mockApiService.EXPECT().GetAllApis().Return(nil, nil)

		mockFactory := NewMockFactory(ctrl)
		mockFactory.EXPECT().Apis().Return(mockApiService).AnyTimes()

		mockClient := &clientv1.Client{
			Factory: mockFactory,
		}

		a := &RouteRegistry{
			group:      group,
			client:     mockClient,
			apiGroup:   "user",
			mu:         sync.Mutex{},
			apiCache:   make(map[string]*model.SysApi),
			cacheOnce:  sync.Once{},
			strictMode: false,
		}

		assert.NotPanics(t, func() {
			a.ensureCacheLoaded()
		})
		assert.Empty(t, a.apiCache)
	})
}

func TestRouteRegistry_ensureCacheLoaded_WithEmptyApis(t *testing.T) {
	t.Run("ensure cache loaded with empty API list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gin.SetMode(gin.TestMode)
		engine := gin.New()
		group := engine.Group("/api")

		mockApiService := service.NewMockApiService(ctrl)
		mockApiService.EXPECT().GetAllApis().Return(&response.SysAPIListResponse{
			Apis: nil,
		}, nil)

		mockFactory := NewMockFactory(ctrl)
		mockFactory.EXPECT().Apis().Return(mockApiService).AnyTimes()

		mockClient := &clientv1.Client{
			Factory: mockFactory,
		}

		a := &RouteRegistry{
			group:      group,
			client:     mockClient,
			apiGroup:   "user",
			mu:         sync.Mutex{},
			apiCache:   make(map[string]*model.SysApi),
			cacheOnce:  sync.Once{},
			strictMode: false,
		}

		assert.NotPanics(t, func() {
			a.ensureCacheLoaded()
		})
		assert.Empty(t, a.apiCache)
	})
}

func TestRouteRegistry_ensureCacheLoaded_WithError(t *testing.T) {
	t.Run("ensure cache loaded with error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gin.SetMode(gin.TestMode)
		engine := gin.New()
		group := engine.Group("/api")

		mockApiService := service.NewMockApiService(ctrl)
		mockApiService.EXPECT().GetAllApis().Return(nil, assert.AnError)

		mockFactory := NewMockFactory(ctrl)
		mockFactory.EXPECT().Apis().Return(mockApiService).AnyTimes()

		mockClient := &clientv1.Client{
			Factory: mockFactory,
		}

		a := &RouteRegistry{
			group:      group,
			client:     mockClient,
			apiGroup:   "user",
			mu:         sync.Mutex{},
			apiCache:   make(map[string]*model.SysApi),
			cacheOnce:  sync.Once{},
			strictMode: false,
		}

		assert.NotPanics(t, func() {
			a.ensureCacheLoaded()
		})
		assert.Empty(t, a.apiCache)
	})
}

func TestRouteRegistry_buildCacheKey(t *testing.T) {
	type fields struct {
		group      *gin.RouterGroup
		client     *clientv1.Client
		apiGroup   string
		mu         sync.Mutex
		apiCache   map[string]*model.SysApi
		cacheOnce  sync.Once
		strictMode bool
	}
	type args struct {
		path     string
		apiGroup string
		method   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "build cache key for GET user list",
			fields: fields{
				apiGroup: "user",
			},
			args: args{
				path:     "/api/v1/users",
				apiGroup: "user",
				method:   "GET",
			},
			want: "GET:user:/api/v1/users",
		},
		{
			name: "build cache key for POST user creation",
			fields: fields{
				apiGroup: "user",
			},
			args: args{
				path:     "/api/v1/users",
				apiGroup: "user",
				method:   "POST",
			},
			want: "POST:user:/api/v1/users",
		},
		{
			name: "build cache key for DELETE operation",
			fields: fields{
				apiGroup: "admin",
			},
			args: args{
				path:     "/api/v1/users/:id",
				apiGroup: "admin",
				method:   "DELETE",
			},
			want: "DELETE:admin:/api/v1/users/:id",
		},
		{
			name: "build cache key with ANY method",
			fields: fields{
				apiGroup: "health",
			},
			args: args{
				path:     "/health",
				apiGroup: "health",
				method:   "ANY",
			},
			want: "ANY:health:/health",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &RouteRegistry{
				group:      tt.fields.group,
				client:     tt.fields.client,
				apiGroup:   tt.fields.apiGroup,
				mu:         tt.fields.mu,
				apiCache:   tt.fields.apiCache,
				cacheOnce:  tt.fields.cacheOnce,
				strictMode: tt.fields.strictMode,
			}
			if got := a.buildCacheKey(tt.args.path, tt.args.apiGroup, tt.args.method); got != tt.want {
				t.Errorf("buildCacheKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteRegistry_ensureCacheLoaded(t *testing.T) {
	type fields struct {
		client     *clientv1.Client
		apiGroup   string
		mu         sync.Mutex
		apiCache   map[string]*model.SysApi
		cacheOnce  sync.Once
		strictMode bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "ensure cache loaded with nil client should panic",
			fields: fields{
				client:     nil,
				apiGroup:   "test",
				apiCache:   make(map[string]*model.SysApi),
				strictMode: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			engine := gin.New()
			group := engine.Group("/api")

			a := &RouteRegistry{
				group:      group,
				client:     tt.fields.client,
				apiGroup:   tt.fields.apiGroup,
				mu:         tt.fields.mu,
				apiCache:   tt.fields.apiCache,
				cacheOnce:  tt.fields.cacheOnce,
				strictMode: tt.fields.strictMode,
			}
			assert.Panics(t, func() {
				a.ensureCacheLoaded()
			})
		})
	}
}

func Test_joinPaths(t *testing.T) {
	type args struct {
		absolutePath string
		relativePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "join base path with relative path",
			args: args{
				absolutePath: "/api/v1",
				relativePath: "/users",
			},
			want: "/api/v1/users",
		},
		{
			name: "join base path with empty relative path",
			args: args{
				absolutePath: "/api/v1",
				relativePath: "",
			},
			want: "/api/v1",
		},
		{
			name: "join base path with relative path ending with slash",
			args: args{
				absolutePath: "/api/v1",
				relativePath: "/users/",
			},
			want: "/api/v1/users/",
		},
		{
			name: "join root path with relative path",
			args: args{
				absolutePath: "/",
				relativePath: "users",
			},
			want: "/users",
		},
		{
			name: "join paths with multiple segments",
			args: args{
				absolutePath: "/api/v1/admin",
				relativePath: "/users/:id",
			},
			want: "/api/v1/admin/users/:id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinPaths(tt.args.absolutePath, tt.args.relativePath); got != tt.want {
				t.Errorf("joinPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lastChar(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		{
			name: "get last char of normal string",
			args: args{
				str: "/api/v1/users",
			},
			want: 's',
		},
		{
			name: "get last char of string ending with slash",
			args: args{
				str: "/api/v1/users/",
			},
			want: '/',
		},
		{
			name: "get last char of single char string",
			args: args{
				str: "/",
			},
			want: '/',
		},
		{
			name: "get last char of path with query params",
			args: args{
				str: "/api/v1/users?id=1",
			},
			want: '1',
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lastChar(tt.args.str); got != tt.want {
				t.Errorf("lastChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lastChar_PanicsOnEmptyString(t *testing.T) {
	assert.Panics(t, func() {
		lastChar("")
	})
}

func TestRouteRegistry_Integration_AutoRegistrationFlow(t *testing.T) {
	t.Run("full auto registration flow mimicking InitHeimdallrApiWithAutoApiRegistration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gin.SetMode(gin.TestMode)
		engine := gin.New()
		rg := engine.Group("/api/clientv1")

		mockApiService := service.NewMockApiService(ctrl)
		mockApiService.EXPECT().GetAllApis().Return(&response.SysAPIListResponse{}, nil).AnyTimes()
		mockApiService.EXPECT().CreateApi(gomock.Any()).Return(nil).AnyTimes()

		mockAuthzService := service.NewMockAuthService(ctrl)
		mockAuthzService.EXPECT().ValidateToken(gomock.Any()).Return(&response.ValidateTokenResponse{
			Claims: &request.CustomClaims{
				ID:          1,
				AuthorityId: "888",
				BufferTime:  86400,
			},
		}, nil).AnyTimes()
		mockAuthzService.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&response.CheckPermissionResponse{
			Allowed: true,
		}, nil).AnyTimes()
		mockAuthzService.EXPECT().LogOperation(gomock.Any()).Return(nil).AnyTimes()

		mockFactory := NewMockFactory(ctrl)
		mockFactory.EXPECT().Apis().Return(mockApiService).AnyTimes()
		mockFactory.EXPECT().Authz().Return(mockAuthzService).AnyTimes()

		mockClient := &clientv1.Client{
			Factory: mockFactory,
		}

		gm, err := New(mockClient, "heimdallr")
		assert.NoError(t, err)
		assert.NotNil(t, gm)

		groupRoutes := gm.Group(rg, "hmdrGroup").
			Use(gm.JWTAuth()).
			Use(gm.CasbinHandler()).
			Use(gm.OperationRecord())
		{
			groupRoutes.
				POST("createHmdrGroup", "新建HmdrGroup").
				DELETE("deleteHmdrGroup", "删除HmdrGroup").
				DELETE("deleteHmdrGroupByIds", "批量删除HmdrGroup").
				PUT("updateHmdrGroup", "更新HmdrGroup").
				GET("findHmdrGroup", "根据ID获取HmdrGroup").
				GET("getHmdrGroupList", "获取HmdrGroup列表")
		}

		hostRoutes := gm.Group(rg, "hmdrHost").
			Use(gm.JWTAuth()).
			Use(gm.CasbinHandler()).
			Use(gm.OperationRecord())
		{
			hostRoutes.
				POST("createHmdrHost", "新建HmdrHost").
				DELETE("deleteHmdrHost", "删除HmdrHost").
				DELETE("deleteHmdrHostByIds", "批量删除HmdrHost").
				PUT("updateHmdrHost", "更新HmdrHost").
				GET("findHmdrHost", "根据ID获取HmdrHost").
				GET("getHmdrHostList", "获取HmdrHost列表")
		}

		webSrvConfRoutes := gm.Group(rg, "conf").
			Use(gm.JWTAuth()).
			Use(gm.CasbinHandler()).
			Use(gm.OperationRecord())
		{
			webSrvConfRoutes.
				GET("getOptions", "获取options选择参数信息").
				POST("getConfInfo", "获取配置文件信息").
				POST("get-context-text", "获取上下文配置明细").
				POST("get-conf-struct", "获取配置文件JSON数据").
				POST("get-includes", "获取包含的配置文件路径列表").
				POST("search-ctx-poses", "搜索上下文坐标列表").
				POST("insert-clone-ctx", "插入克隆上下文").
				POST("insert-new-ctx", "插入新上下文").
				DELETE("remove-ctx", "移除上下文").
				POST("update-conf", "更新配置").
				POST("modify-ctx-value", "修改上下文值").
				POST("modify-clone-ctx", "修改克隆上下文").
				POST("change-ctx-enabled-state", "修改上下文启用状态").
				POST("modify-new-ctx", "修改新上下文").
				POST("move-ctx", "移动上下文")
		}

		webSrvBinCMDRoutes := gm.Group(rg, "bin-cmd").
			Use(gm.JWTAuth()).
			Use(gm.CasbinHandler()).
			Use(gm.OperationRecord())
		{
			webSrvBinCMDRoutes.POST("exec", "提交Web服务端二进制命令工具执行请求")
		}

		webSrvStatisticsRoutes := gm.Group(rg, "hmdr-statistics").
			Use(gm.JWTAuth()).
			Use(gm.CasbinHandler()).
			Use(gm.OperationRecord())
		{
			webSrvStatisticsRoutes.
				POST("proxy-svc-brief", "获取代理配置信息").
				POST("conn-check-of-proxy-svc", "代理服务网络连通性检查").
				POST("export-proxy-svc-excel", "导出代理信息为Excel")
		}

		assert.NotPanics(t, func() {
			_ = groupRoutes
			_ = hostRoutes
			_ = webSrvConfRoutes
			_ = webSrvBinCMDRoutes
			_ = webSrvStatisticsRoutes
		})
	})
}

func TestRouteRegistry_Integration_MinimalUsage(t *testing.T) {
	t.Run("minimal usage example with CRUD operations", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gin.SetMode(gin.TestMode)
		engine := gin.New()
		rg := engine.Group("/api/clientv1")

		mockApiService := service.NewMockApiService(ctrl)
		mockApiService.EXPECT().GetAllApis().Return(&response.SysAPIListResponse{}, nil).AnyTimes()
		mockApiService.EXPECT().CreateApi(gomock.Any()).Return(nil).AnyTimes()

		mockFactory := NewMockFactory(ctrl)
		mockFactory.EXPECT().Apis().Return(mockApiService).AnyTimes()

		mockClient := &clientv1.Client{
			Factory: mockFactory,
		}

		gm, err := New(mockClient, "example")
		assert.NoError(t, err)
		assert.NotNil(t, gm)

		dummyHandler := func(c *gin.Context) { c.String(200, "ok") }

		groupRoutes := gm.Group(rg, "groups")
		{
			groupRoutes.
				POST("create", "创建资源", dummyHandler).
				GET("list", "获取资源列表", dummyHandler).
				PUT("update", "更新资源", dummyHandler).
				DELETE("delete", "删除资源", dummyHandler)
		}

		assert.NotPanics(t, func() {
			_ = groupRoutes
		})
	})
}
