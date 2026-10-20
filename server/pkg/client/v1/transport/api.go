package transport

import (
	"github.com/tanganyu1114/heimdallr-reborn/server/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	http_transport "github.com/go-kit/kit/transport/http"
)

// ApiTransport defines the interface for API management related transport
type ApiTransport interface {
	// CreateApi returns the create API client
	CreateApi() httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	// UpdateApi returns the update API client
	UpdateApi() httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	// GetAllApis returns the get all APIs client
	GetAllApis() httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]]
}

// apiTransport implements ApiTransport interface
type apiTransport struct {
	createApiClient  httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	updateApiClient  httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	getAllApisClient httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]]
}

// newApiTransport creates a new API transport
func newApiTransport(transport *transport) ApiTransport {
	t := &apiTransport{
		createApiClient: httpclientv1.NewClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/api/createApi",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		updateApiClient: httpclientv1.NewClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/api/updateApi",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
		getAllApisClient: httpclientv1.NewClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]](
			httpclientv1.HTTPMethodPost,
			transport.baseURL+"/api/getAllApis",
		).WithOptions(
			http_transport.SetClient(transport.Client),
		),
	}
	return t
}

// CreateApi returns the create API client
func (a *apiTransport) CreateApi() httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	return a.createApiClient
}

// UpdateApi returns the update API client
func (a *apiTransport) UpdateApi() httpclientv1.ClientBuilder[*model.SysApi, modelclientv1.ResponseBody[httpclientv1.NilBody]] {
	return a.updateApiClient
}

// GetAllApis returns the get all APIs client
func (a *apiTransport) GetAllApis() httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[*response.SysAPIListResponse]] {
	return a.getAllApisClient
}
