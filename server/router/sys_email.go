package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/heimdallr-reborn/server/api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/middleware"
)

func InitEmailRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("email").Use(middleware.OperationRecord())
	{
		UserRouter.POST("emailTest", v1.EmailTest) // 发送测试邮件
	}
}
