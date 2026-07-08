package web_server_statistics

import (
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/service/v1"
)

type WebServerStatisticsController struct {
	svc svcv1.Factory
}

func NewController(service svcv1.Factory) *WebServerStatisticsController {
	return &WebServerStatisticsController{
		svc: service,
	}
}
