package router

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

func InitHmdrConfRouter(Router *gin.RouterGroup) {
	InitHmdrConfRouter := Router.Group("conf").Use(middleware.OperationRecord())
	{
		InitHmdrConfRouter.GET("getOptions", v1.GetOptions)    // 获取options选择参数信息
		InitHmdrConfRouter.POST("getConfInfo", v1.GetConfInfo) // 获取配置文件信息
	}
}
