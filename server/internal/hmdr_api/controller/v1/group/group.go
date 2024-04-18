package group

import (
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
)

type GroupController struct {
	svc svcv1.Factory
}

func NewController(service svcv1.Factory) *GroupController {
	return &GroupController{
		svc: service,
	}
}
