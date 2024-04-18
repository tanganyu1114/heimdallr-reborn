package host

import (
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
)

type HostController struct {
	svc svcv1.Factory
}

func NewController(service svcv1.Factory) *HostController {
	return &HostController{
		svc: service,
	}
}
