package service

import (
	"context"
	"github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
	epclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/endpoint"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/marmotedu/errors"
)

// HostService defines the interface for host related services
type HostService interface {
	// Get retrieves a host by ID
	Get(idOptions *metav1.IDOptions) (*v1.Host, error)
	// List retrieves a list of hosts
	List(listOptions *metav1.ListOptions) (*v1.HostList, error)
}

// hostService implements HostService interface
type hostService struct {
	ctx context.Context
	eps epclientv1.HostEndpoints
}

// newHostService creates a new Host service
func newHostService(svc *factory) HostService {
	if svc == nil || svc.eps == nil {
		return nil
	}
	return &hostService{
		ctx: svc.ctx,
		eps: svc.eps.Hosts(),
	}
}

// Get retrieves a host by ID
func (s *hostService) Get(idOptions *metav1.IDOptions) (*v1.Host, error) {
	if idOptions == nil {
		return nil, errors.New("idOptions cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.IDOptions]{
		Body: *idOptions,
	}
	resp, err := s.eps.Get()(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Response()
}

// List retrieves a list of hosts
func (s *hostService) List(listOptions *metav1.ListOptions) (*v1.HostList, error) {
	if listOptions == nil {
		return nil, errors.New("listOptions cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.ListOptions]{
		Body: *listOptions,
	}
	resp, err := s.eps.List()(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Response()
}
