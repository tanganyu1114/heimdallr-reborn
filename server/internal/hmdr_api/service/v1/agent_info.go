package v1

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
)

type AgentInfoSrv interface {
	Get(ctx context.Context) ([]v1.GroupInfo, error)
}

type agentInfoService struct {
	store storev1.Factory
}

var _ AgentInfoSrv = (*agentInfoService)(nil)

func (a *agentInfoService) Get(ctx context.Context) ([]v1.GroupInfo, error) {
	return a.store.AgentInfos().Get(ctx)
}

func newAgentInfos(svc *service) AgentInfoSrv {
	return &agentInfoService{
		store: svc.store,
	}
}
