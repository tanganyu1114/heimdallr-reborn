package service

import (
	"context"
	"gin-vue-admin/api/heimdallr_api/v1"
	epclientv1 "gin-vue-admin/pkg/client/v1/endpoint"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

// AgentInfoService defines the interface for agent info related services
type AgentInfoService interface {
	// Get retrieves agent information
	Get() ([]v1.GroupInfo, error)
}

// agentInfoService implements AgentInfoService interface
type agentInfoService struct {
	ctx context.Context
	eps epclientv1.AgentInfoEndpoints
}

// newAgentInfoService creates a new AgentInfo service
func newAgentInfoService(svc *factory) AgentInfoService {
	if svc == nil || svc.eps == nil {
		return nil
	}
	return &agentInfoService{
		ctx: svc.ctx,
		eps: svc.eps.AgentInfos(),
	}
}

// Get retrieves agent information
func (s *agentInfoService) Get() ([]v1.GroupInfo, error) {
	req := httpclientv1.HTTPRequest[httpclientv1.NilBody]{}
	return s.eps.Get()(s.ctx, req)
}
