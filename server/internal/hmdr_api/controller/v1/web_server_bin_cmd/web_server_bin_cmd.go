package web_server_bin_cmd

import (
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/service/v1"
)

type WebServerBinCMDController struct {
	svc svcv1.Factory
}

func NewController(service svcv1.Factory) *WebServerBinCMDController {
	return &WebServerBinCMDController{
		svc: service,
	}
}
