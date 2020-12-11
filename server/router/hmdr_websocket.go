package router

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

func InitHmdrWebSocketRouter(Router *gin.RouterGroup) {
	HmdrWebSocketRouter := Router.Group("hmdrWebSocket").Use(middleware.OperationRecord())
	{
		HmdrWebSocketRouter.GET("ws", v1.WebSocket)
	}
}
