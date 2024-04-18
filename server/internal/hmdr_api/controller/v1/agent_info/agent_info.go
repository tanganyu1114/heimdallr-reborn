package agent_info

import (
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
)

type AgentInfoController struct {
	svc svcv1.Factory
}

func NewController(service svcv1.Factory) *AgentInfoController {
	return &AgentInfoController{
		svc: service,
	}
}
