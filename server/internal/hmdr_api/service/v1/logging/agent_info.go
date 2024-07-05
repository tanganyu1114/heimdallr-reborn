package logging

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	"go.uber.org/zap/zapcore"
)

type agentInfoService struct {
	svc svcv1.Factory
}

var _ svcv1.AgentInfoSrv = (*agentInfoService)(nil)

func (a *agentInfoService) Get(ctx context.Context) (groupInfos []v1.GroupInfo, err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "查询服务状态", nil, groupInfos, err)
	}()
	return a.svc.AgentInfos().Get(ctx)
}

func newAgentInfos(svc svcv1.Factory) *agentInfoService {
	return &agentInfoService{svc: svc}
}
