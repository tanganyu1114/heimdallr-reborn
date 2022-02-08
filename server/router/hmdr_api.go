package router

import (
	"gin-vue-admin/internal/hmdr_api/controller/v1/agent_info"
	"gin-vue-admin/internal/hmdr_api/controller/v1/group"
	"gin-vue-admin/internal/hmdr_api/controller/v1/host"
	"gin-vue-admin/internal/hmdr_api/controller/v1/web_server_config"
	"gin-vue-admin/internal/hmdr_api/controller/v1/web_server_log_watcher"
	"gin-vue-admin/internal/hmdr_api/store/v1/bifrosts"
	"gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

func InitPublicHeimdallrApi(rg *gin.RouterGroup) {
	storeIns := bifrosts.GetBifrostsStore()
	webSrvLogWatcherRoutes := rg.Group("hmdrWebSocket").Use(middleware.OperationRecord())
	{
		webSrvLogWatcherController := web_server_log_watcher.NewController(storeIns)

		webSrvLogWatcherRoutes.GET("ws", webSrvLogWatcherController.Watch)
	}
}

func InitPrivateHeimdallrApi(rg *gin.RouterGroup) {
	storeIns := bifrosts.GetBifrostsStore()
	agentRoutes := rg.Group("agent").Use(middleware.OperationRecord())
	{
		agentController := agent_info.NewController(storeIns)

		agentRoutes.GET("getAgentInfo", agentController.Get)
	}

	groupRoutes := rg.Group("hmdrGroup").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler()).Use(middleware.OperationRecord())
	{
		groupController := group.NewController(storeIns)

		groupRoutes.POST("createHmdrGroup", groupController.Create)                   // 新建HmdrGroup
		groupRoutes.DELETE("deleteHmdrGroup", groupController.Delete)                 // 删除HmdrGroup
		groupRoutes.DELETE("deleteHmdrGroupByIds", groupController.DeleteCollections) // 批量删除HmdrGroup
		groupRoutes.PUT("updateHmdrGroup", groupController.Update)                    // 更新HmdrGroup
		groupRoutes.GET("findHmdrGroup", groupController.Get)                         // 根据ID获取HmdrGroup
		groupRoutes.GET("getHmdrGroupList", groupController.List)                     // 获取HmdrGroup列表
	}

	hostRoutes := rg.Group("hmdrHost").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler()).Use(middleware.OperationRecord())
	{
		hostController := host.NewController(storeIns)

		hostRoutes.POST("createHmdrHost", hostController.Create)                  // 新建HmdrHost
		hostRoutes.DELETE("deleteHmdrHost", hostController.Delete)                // 删除HmdrHost
		hostRoutes.DELETE("deleteHmdrHostByIds", hostController.DeleteCollection) // 批量删除HmdrHost
		hostRoutes.PUT("updateHmdrHost", hostController.Update)                   // 更新HmdrHost
		hostRoutes.GET("findHmdrHost", hostController.Get)                        // 根据ID获取HmdrHost
		hostRoutes.GET("getHmdrHostList", hostController.List)                    // 获取HmdrHost列表
	}

	webSrvConfRoutes := rg.Group("conf").Use(middleware.OperationRecord())
	{
		webSrvConfController := web_server_config.NewController(storeIns)

		webSrvConfRoutes.GET("getOptions", webSrvConfController.GetOptions)  // 获取options选择参数信息
		webSrvConfRoutes.POST("getConfInfo", webSrvConfController.GetConfig) // 获取配置文件信息
	}
}
