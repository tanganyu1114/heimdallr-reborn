package router

import (
	"github.com/tanganyu1114/heimdallr-reborn/api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/middleware"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.POST("login", v1.Login)
		BaseRouter.POST("sdkLogin", middleware.OperationRecord(), v1.SDKLogin)
		BaseRouter.POST("captcha", v1.Captcha)
	}
	return BaseRouter
}
