package web_server_statistics

import (
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
)

type WebServerStatisticsController struct {
	svc svcv1.Service
}

func NewController(store storev1.Factory) *WebServerStatisticsController {
	return &WebServerStatisticsController{
		svc: svcv1.NewService(store),
	}
}
