package bifrosts

import (
	"context"
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/service/v1"
	storev1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1"
)

type agentInfoService struct {
	store storev1.Factory
}

var _ svcv1.AgentInfoSrv = (*agentInfoService)(nil)

func (a *agentInfoService) Get(ctx context.Context) ([]v1.GroupInfo, error) {
	return a.store.AgentInfos().Get(ctx)
}

func newAgentInfos(svc *service) svcv1.AgentInfoSrv {
	return &agentInfoService{
		store: svc.store,
	}
}
