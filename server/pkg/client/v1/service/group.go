package service

import (
	"context"
	"gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	epclientv1 "gin-vue-admin/pkg/client/v1/endpoint"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/marmotedu/errors"
)

// GroupService defines the interface for group related services
type GroupService interface {
	// Get retrieves a group by ID
	Get(idOptions *metav1.IDOptions) (*v1.Group, error)
	// List retrieves a list of groups
	List(listOptions *metav1.ListOptions) (*v1.GroupList, error)
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
func (s *groupService) Get(idOptions *metav1.IDOptions) (*v1.Group, error) {
	if idOptions == nil {
		return nil, errors.New("idOptions cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.IDOptions]{
		Body: *idOptions,
	}
	return s.eps.Get()(s.ctx, req)
}

// List retrieves a list of groups
func (s *groupService) List(listOptions *metav1.ListOptions) (*v1.GroupList, error) {
	if listOptions == nil {
		return nil, errors.New("listOptions cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.ListOptions]{
		Body: *listOptions,
	}
	return s.eps.List()(s.ctx, req)
}
