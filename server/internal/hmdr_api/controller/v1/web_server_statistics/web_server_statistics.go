package web_server_statistics

import (
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
)

type WebServerStatisticsController struct {
	svc svcv1.Factory
}

func NewController(service svcv1.Factory) *WebServerStatisticsController {
	return &WebServerStatisticsController{
		svc: service,
	}
}
