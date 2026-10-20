package ginmiddleware

import (
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/ClessLi/component-base/pkg/errors"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/heimdallr-reborn/server/model"
	clientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1"
)

// RouteRegistry is a decorator for gin.RouterGroup that automatically creates API entries
// when routes are registered
type RouteRegistry struct {
	group      *gin.RouterGroup
	client     *clientv1.Client
	apiGroup   string
	mu         sync.Mutex
	apiCache   map[string]*model.SysApi
	cacheOnce  sync.Once
	strictMode bool
}

func (gm *GinMiddleware) Group(rg *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) *RouteRegistry {
	if rg == nil {
		panic("RouterGroup cannot be nil")
	}
	return (&RouteRegistry{
		group:    rg,
		client:   gm.client,
		apiGroup: gm.apiGroup,
	}).Group(relativePath, handlers...)
}

// Group creates a new router group with the given relative path and middleware
func (a *RouteRegistry) Group(relativePath string, handlers ...gin.HandlerFunc) *RouteRegistry {
	return a.registerRoute(relativePath, "", "", handlers)
}

// Use adds middleware to the router group
func (a *RouteRegistry) Use(middleware ...gin.HandlerFunc) *RouteRegistry {
	a.group.Use(middleware...)
	return a
}

// POST registers a POST route and automatically creates an API entry
func (a *RouteRegistry) POST(relativePath, desc string, handlers ...gin.HandlerFunc) *RouteRegistry {
	return a.registerRoute(relativePath, http.MethodPost, desc, handlers)
}

// GET registers a GET route and automatically creates an API entry
func (a *RouteRegistry) GET(relativePath, desc string, handlers ...gin.HandlerFunc) *RouteRegistry {
	return a.registerRoute(relativePath, http.MethodGet, desc, handlers)
}

// PUT registers a PUT route and automatically creates an API entry
func (a *RouteRegistry) PUT(relativePath, desc string, handlers ...gin.HandlerFunc) *RouteRegistry {
	return a.registerRoute(relativePath, http.MethodPut, desc, handlers)
}

// DELETE registers a DELETE route and automatically creates an API entry
func (a *RouteRegistry) DELETE(relativePath, desc string, handlers ...gin.HandlerFunc) *RouteRegistry {
	return a.registerRoute(relativePath, http.MethodDelete, desc, handlers)
}

// PATCH registers a PATCH route and automatically creates an API entry
func (a *RouteRegistry) PATCH(relativePath, desc string, handlers ...gin.HandlerFunc) *RouteRegistry {
	return a.registerRoute(relativePath, http.MethodPatch, desc, handlers)
}

// HEAD registers a HEAD route and automatically creates an API entry
func (a *RouteRegistry) HEAD(relativePath, desc string, handlers ...gin.HandlerFunc) *RouteRegistry {
	return a.registerRoute(relativePath, http.MethodHead, desc, handlers)
}

// OPTIONS registers an OPTIONS route and automatically creates an API entry
func (a *RouteRegistry) OPTIONS(relativePath, desc string, handlers ...gin.HandlerFunc) *RouteRegistry {
	return a.registerRoute(relativePath, http.MethodOptions, desc, handlers)
}

// Any registers a route for all HTTP methods and automatically creates an API entry
func (a *RouteRegistry) Any(relativePath, desc string, handlers ...gin.HandlerFunc) *RouteRegistry {
	return a.registerRoute(relativePath, "ANY", desc, handlers)
}

// registerRoute registers a route with the given method, description, and handlers
// and automatically creates an API entry
func (a *RouteRegistry) registerRoute(relativePath, method, desc string, handlers []gin.HandlerFunc) *RouteRegistry {
	if method != "" {
		a.applyAPIEntry(method, relativePath, desc)
		a.routerHandle(method, relativePath, handlers)
	} else {
		a.group = a.group.Group(relativePath, handlers...)
	}
	return &RouteRegistry{
		group:    a.group,
		client:   a.client,
		apiGroup: a.apiGroup,
	}
}

// routerHandle registers the actual route with gin based on the HTTP method.
// It dispatches the route registration to the appropriate gin.RouterGroup method
// (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, or Any) according to the method parameter.
func (a *RouteRegistry) routerHandle(method, relativePath string, handlers []gin.HandlerFunc) {
	switch method {
	case http.MethodGet:
		a.group.GET(relativePath, handlers...)
	case http.MethodPost:
		a.group.POST(relativePath, handlers...)
	case http.MethodPut:
		a.group.PUT(relativePath, handlers...)
	case http.MethodDelete:
		a.group.DELETE(relativePath, handlers...)
	case http.MethodPatch:
		a.group.PATCH(relativePath, handlers...)
	case http.MethodHead:
		a.group.HEAD(relativePath, handlers...)
	case http.MethodOptions:
		a.group.OPTIONS(relativePath, handlers...)
	case "ANY":
		a.group.Any(relativePath, handlers...)
	}
}

// applyAPIEntry creates an API entry via SDK service
func (a *RouteRegistry) applyAPIEntry(method, subPath, description string) {
	if a.client == nil {
		return
	}

	// Build full API subPath
	fullPath := joinPaths(a.group.BasePath(), subPath)

	// Create API entry
	api := &model.SysApi{
		Path:        fullPath,
		Description: description,
		ApiGroup:    a.apiGroup,
		Method:      strings.ToUpper(method),
	}

	err := a.applyAPI(api)
	if err != nil {
		if a.strictMode {
			panic("Failed to register API entry: " + err.Error())
		}
		logV1.Errorf("Failed to apply API entry: %+v", err)
	}
}

// applyAPI creates or updates an API entry in the system.
// It validates the API object, loads existing API cache, and either updates
// an existing entry or creates a new one based on the cache lookup.
// Returns error if validation fails or API operations encounter issues.
func (a *RouteRegistry) applyAPI(api *model.SysApi) error {
	if api == nil {
		return errors.New("Ensure API object is initialized before registration")
	}
	if api.Path == "" {
		return errors.New("Provide a valid API path for route registration")
	}
	if api.ApiGroup == "" {
		return errors.New("Specify API group to categorize the route properly")
	}
	if api.Method == "" {
		return errors.New("Define HTTP method for the API route")
	}
	if api.Description == "" {
		return errors.New("Add description to document the API endpoint purpose")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.ensureCacheLoaded()

	cacheKey := a.buildCacheKey(api.Path, api.ApiGroup, api.Method)
	if existingAPI, exists := a.apiCache[cacheKey]; exists {
		if api.Description != existingAPI.Description {
			api.GVA_MODEL = existingAPI.GVA_MODEL
			if err := a.client.Apis().UpdateApi(api); err != nil {
				return errors.Wrap(err, "Check API update permissions and retry")
			}
			a.apiCache[cacheKey] = api
		}
		return nil
	}

	if err := a.client.Apis().CreateApi(api); err != nil {
		return errors.Wrap(err, "Verify API registration permissions and retry")
	}

	a.apiCache[cacheKey] = api
	return nil
}

func (a *RouteRegistry) ensureCacheLoaded() {
	a.cacheOnce.Do(func() {
		a.apiCache = make(map[string]*model.SysApi)
		resp, err := a.client.Apis().GetAllApis()
		if err != nil {
			logV1.Errorf("Failed to load API cache: %+v", err)
			return
		}
		if resp == nil || resp.Apis == nil {
			return
		}
		for i := range resp.Apis {
			existingAPI := &resp.Apis[i]
			key := a.buildCacheKey(existingAPI.Path, existingAPI.ApiGroup, existingAPI.Method)
			a.apiCache[key] = existingAPI
		}
	})
}

func (a *RouteRegistry) buildCacheKey(path, apiGroup, method string) string {
	return method + ":" + apiGroup + ":" + path
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}
