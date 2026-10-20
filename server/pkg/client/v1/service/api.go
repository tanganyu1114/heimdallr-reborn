package service

import (
	"context"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	epclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/endpoint"
)

// ApiService defines the interface for API management service
type ApiService interface {
	// CreateApi creates a new API entry
	CreateApi(api *model.SysApi) error
	// UpdateApi updates an existing API entry
	UpdateApi(api *model.SysApi) error
	// GetAllApis retrieves all registered APIs
	GetAllApis() (*response.SysAPIListResponse, error)
}

// apiService implements ApiService interface
type apiService struct {
	ctx context.Context
	eps epclientv1.ApiEndpoints
}

// newApiService creates a new API service
func newApiService(svc *factory) ApiService {
	if svc == nil || svc.eps == nil {
		return nil
	}
	return &apiService{
		ctx: svc.ctx,
		eps: svc.eps.Apis(),
	}
}

// CreateApi creates a new API entry
func (a *apiService) CreateApi(api *model.SysApi) error {
	req := httpclientv1.HTTPRequest[*model.SysApi]{
		Body: api,
	}
	resp, err := a.eps.CreateApi()(a.ctx, req)
	if err != nil {
		return err
	}
	_, err = resp.Response()
	return err
}

// UpdateApi updates an existing API entry
func (a *apiService) UpdateApi(api *model.SysApi) error {
	req := httpclientv1.HTTPRequest[*model.SysApi]{
		Body: api,
	}
	resp, err := a.eps.UpdateApi()(a.ctx, req)
	if err != nil {
		return err
	}
	_, err = resp.Response()
	return err
}

// GetAllApis retrieves all registered APIs
func (a *apiService) GetAllApis() (*response.SysAPIListResponse, error) {
	req := httpclientv1.HTTPRequest[httpclientv1.NilBody]{
		Body: httpclientv1.NilBody{},
	}
	resp, err := a.eps.GetAllApis()(a.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Response()
}
