package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/heimdallr-reborn/api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/middleware"
)

func InitSystemRouter(Router *gin.RouterGroup) {
	SystemRouter := Router.Group("system").Use(middleware.OperationRecord())
	{
		SystemRouter.POST("getSystemConfig", v1.GetSystemConfig) // 获取配置文件内容
		SystemRouter.POST("setSystemConfig", v1.SetSystemConfig) // 设置配置文件内容
		SystemRouter.POST("getServerInfo", v1.GetServerInfo)     // 获取服务器信息
	}
}
