package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/heimdallr-reborn/api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/middleware"
)

func InitCasbinRouter(Router *gin.RouterGroup) {
	CasbinRouter := Router.Group("casbin").Use(middleware.OperationRecord())
	{
		CasbinRouter.POST("updateCasbin", v1.UpdateCasbin)
		CasbinRouter.POST("getPolicyPathByAuthorityId", v1.GetPolicyPathByAuthorityId)
	}
}
