package router

import (
	"gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

func InitHmdrHostRouter(Router *gin.RouterGroup) {
	HmdrHostRouter := Router.Group("hmdrHost").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler()).Use(middleware.OperationRecord())
	{
		HmdrHostRouter.POST("createHmdrHost", v1.CreateHmdrHost)             // 新建HmdrHost
		HmdrHostRouter.DELETE("deleteHmdrHost", v1.DeleteHmdrHost)           // 删除HmdrHost
		HmdrHostRouter.DELETE("deleteHmdrHostByIds", v1.DeleteHmdrHostByIds) // 批量删除HmdrHost
		HmdrHostRouter.PUT("updateHmdrHost", v1.UpdateHmdrHost)              // 更新HmdrHost
		HmdrHostRouter.GET("findHmdrHost", v1.FindHmdrHost)                  // 根据ID获取HmdrHost
		HmdrHostRouter.GET("getHmdrHostList", v1.GetHmdrHostList)            // 获取HmdrHost列表
	}
}
