package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/heimdallr-reborn/server/api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/middleware"
)

func InitAuthRouter(Router *gin.RouterGroup) {
	AuthRouter := Router.Group("auth").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler()).Use(middleware.OperationRecord())
	{
		AuthRouter.POST("validateToken", v1.ValidateToken)
		AuthRouter.POST("refreshToken", v1.RefreshToken)
		AuthRouter.POST("checkPermission", v1.CheckPermission)
		AuthRouter.POST("logOperation", v1.LogOperation)
	}
}
