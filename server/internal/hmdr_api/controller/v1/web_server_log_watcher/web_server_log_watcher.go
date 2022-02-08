package web_server_log_watcher

import (
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
)

type WebServerLogWatcherController struct {
	svc svcv1.Service
}

func NewController(store storev1.Factory) *WebServerLogWatcherController {
	return &WebServerLogWatcherController{
		svc: svcv1.NewService(store),
	}
}
