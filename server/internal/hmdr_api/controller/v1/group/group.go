package group

import (
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
)

type GroupController struct {
	svc svcv1.Service
}

func NewController(store storev1.Factory) *GroupController {
	return &GroupController{
		svc: svcv1.NewService(store),
	}
}
