package web_server_log_watcher

import (
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
)

type WebServerLogWatcherController struct {
	svc svcv1.Factory
}

func NewController(service svcv1.Factory) *WebServerLogWatcherController {
	return &WebServerLogWatcherController{
		svc: service,
	}
}
