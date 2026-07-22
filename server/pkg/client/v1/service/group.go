package service

import (
	"context"

	"github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	epclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/endpoint"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/marmotedu/errors"
)

// GroupService defines the interface for group related services
type GroupService interface {
	// Get retrieves a group by ID
	Get(idOptions *v1.IDOptions) (*v1.Group, error)
	// List retrieves a list of groups
	List(listOptions *v1.ListOptions) (*v1.GroupList, error)
}

// groupService implements GroupService interface
type groupService struct {
	ctx context.Context
	eps epclientv1.GroupEndpoints
}

// newGroupService creates a new Group service
func newGroupService(svc *factory) GroupService {
	if svc == nil || svc.eps == nil {
		return nil
	}
	return &groupService{
		ctx: svc.ctx,
		eps: svc.eps.Groups(),
	}
}

// Get retrieves a group by ID
func (s *groupService) Get(idOptions *v1.IDOptions) (*v1.Group, error) {
	if idOptions == nil {
		return nil, errors.New("idOptions cannot be nil")
	}
	req := httpclientv1.HTTPRequest[v1.IDOptions]{
		Body: *idOptions,
	}
	resp, err := s.eps.Get()(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Response()
}

// List retrieves a list of groups
func (s *groupService) List(listOptions *v1.ListOptions) (*v1.GroupList, error) {
	if listOptions == nil {
		return nil, errors.New("listOptions cannot be nil")
	}
	req := httpclientv1.HTTPRequest[v1.ListOptions]{
		Body: *listOptions,
	}
	resp, err := s.eps.List()(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Response()
}
