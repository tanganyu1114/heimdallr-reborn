package router

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

func InitHmdrAgentRouter(Router *gin.RouterGroup) {
	HmdrAgentRouter := Router.Group("agent").Use(middleware.OperationRecord())
	{
		HmdrAgentRouter.GET("getAgentInfo", v1.GetAgentInfo) // 获取biforst agent信息
	}
}
