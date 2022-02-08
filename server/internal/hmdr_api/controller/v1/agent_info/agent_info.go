package agent_info

import (
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
)

type AgentInfoController struct {
	svc svcv1.Service
}

func NewController(store storev1.Factory) *AgentInfoController {
	return &AgentInfoController{
		svc: svcv1.NewService(store),
	}
}
