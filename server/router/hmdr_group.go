package router

import (
	"gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

func InitHmdrGroupRouter(Router *gin.RouterGroup) {
	HmdrGroupRouter := Router.Group("hmdrGroup").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler()).Use(middleware.OperationRecord())
	{
		HmdrGroupRouter.POST("createHmdrGroup", v1.CreateHmdrGroup)             // 新建HmdrGroup
		HmdrGroupRouter.DELETE("deleteHmdrGroup", v1.DeleteHmdrGroup)           // 删除HmdrGroup
		HmdrGroupRouter.DELETE("deleteHmdrGroupByIds", v1.DeleteHmdrGroupByIds) // 批量删除HmdrGroup
		HmdrGroupRouter.PUT("updateHmdrGroup", v1.UpdateHmdrGroup)              // 更新HmdrGroup
		HmdrGroupRouter.GET("findHmdrGroup", v1.FindHmdrGroup)                  // 根据ID获取HmdrGroup
		HmdrGroupRouter.GET("getHmdrGroupList", v1.GetHmdrGroupList)            // 获取HmdrGroup列表
	}
}
