package web_server_config

import (
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
)

type WebServerConfigController struct {
	svc svcv1.Factory
}

func NewController(service svcv1.Factory) *WebServerConfigController {
	return &WebServerConfigController{
		svc: service,
	}
}
