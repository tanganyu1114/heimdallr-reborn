package web_server_config

import (
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
)

type WebServerConfigController struct {
	svc svcv1.Service
}

func NewController(store storev1.Factory) *WebServerConfigController {
	return &WebServerConfigController{
		svc: svcv1.NewService(store),
	}
}
