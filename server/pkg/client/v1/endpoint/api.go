package endpoint

import (
	"sync"

	"github.com/tanganyu1114/heimdallr-reborn/server/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	txpclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/transport"

	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

// ApiEndpoints defines the interface for API management related endpoints
type ApiEndpoints interface {
	// CreateApi returns the create API endpoint
	CreateApi() httpclientv1.Endpoint[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	// UpdateApi returns the update API endpoint
	UpdateApi() httpclientv1.Endpoint[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	// GetAllApis returns the get all APIs endpoint
	GetAllApis() httpclientv1.Endpoint[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]]
}

// apiEndpoints implements ApiEndpoints interface
type apiEndpoints struct {
	transport          txpclientv1.ApiTransport
	onceCreateApi      sync.Once
	createApiEndpoint  httpclientv1.Endpoint[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	onceUpdateApi      sync.Once
	updateApiEndpoint  httpclientv1.Endpoint[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	onceGetAllApis     sync.Once
	getAllApisEndpoint httpclientv1.Endpoint[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]]
}

// CreateApi returns the create API endpoint
func (a *apiEndpoints) CreateApi() httpclientv1.Endpoint[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	a.onceCreateApi.Do(func() {
		a.createApiEndpoint = a.transport.CreateApi().Build().Endpoint()
	})
	return a.createApiEndpoint
}

// UpdateApi returns the update API endpoint
func (a *apiEndpoints) UpdateApi() httpclientv1.Endpoint[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	a.onceUpdateApi.Do(func() {
		a.updateApiEndpoint = a.transport.UpdateApi().Build().Endpoint()
	})
	return a.updateApiEndpoint
}

// GetAllApis returns the get all APIs endpoint
func (a *apiEndpoints) GetAllApis() httpclientv1.Endpoint[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]] {
	a.onceGetAllApis.Do(func() {
		a.getAllApisEndpoint = a.transport.GetAllApis().Build().Endpoint()
	})
	return a.getAllApisEndpoint
}

// newApiEndpoints creates a new API endpoints
func newApiEndpoints(f *factory) ApiEndpoints {
	return &apiEndpoints{
		transport: f.transport.Apis(),
	}
}
