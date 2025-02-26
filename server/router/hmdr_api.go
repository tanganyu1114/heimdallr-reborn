package router

import (
	"gin-vue-admin/internal/hmdr_api/controller/v1/agent_info"
	"gin-vue-admin/internal/hmdr_api/controller/v1/group"
	"gin-vue-admin/internal/hmdr_api/controller/v1/host"
	"gin-vue-admin/internal/hmdr_api/controller/v1/web_server_bin_cmd"
	"gin-vue-admin/internal/hmdr_api/controller/v1/web_server_config"
	"gin-vue-admin/internal/hmdr_api/controller/v1/web_server_log_watcher"
	"gin-vue-admin/internal/hmdr_api/controller/v1/web_server_statistics"
	bifrostssvc "gin-vue-admin/internal/hmdr_api/service/v1/bifrosts"
	loggingsvc "gin-vue-admin/internal/hmdr_api/service/v1/logging"
	bifrostsstore "gin-vue-admin/internal/hmdr_api/store/v1/bifrosts"
	cachestore "gin-vue-admin/internal/hmdr_api/store/v1/cache"
	"gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
	"time"
)

func InitPublicHeimdallrApi(rg *gin.RouterGroup) {
	storeIns := bifrostsstore.GetBifrostsStore()
	svcIns := bifrostssvc.NewService(storeIns)
	webSrvLogWatcherRoutes := rg.Group("hmdrWebSocket").Use(middleware.OperationRecord())
	{
		webSrvLogWatcherController := web_server_log_watcher.NewController(svcIns)

		webSrvLogWatcherRoutes.GET("ws", webSrvLogWatcherController.Watch)
	}
}

func InitPrivateHeimdallrApi(rg *gin.RouterGroup) {
	storeIns := cachestore.GetCacheStore(bifrostsstore.GetBifrostsStore(), time.Minute)
	svcIns := loggingsvc.NewService(bifrostssvc.NewService(storeIns))
	agentRoutes := rg.Group("agent").Use(middleware.OperationRecord())
	{
		agentController := agent_info.NewController(svcIns)

		agentRoutes.GET("getAgentInfo", agentController.Get)
	}

	groupRoutes := rg.Group("hmdrGroup").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler()).Use(middleware.OperationRecord())
	{
		groupController := group.NewController(svcIns)

		groupRoutes.POST("createHmdrGroup", groupController.Create)                   // 新建HmdrGroup
		groupRoutes.DELETE("deleteHmdrGroup", groupController.Delete)                 // 删除HmdrGroup
		groupRoutes.DELETE("deleteHmdrGroupByIds", groupController.DeleteCollections) // 批量删除HmdrGroup
		groupRoutes.PUT("updateHmdrGroup", groupController.Update)                    // 更新HmdrGroup
		groupRoutes.GET("findHmdrGroup", groupController.Get)                         // 根据ID获取HmdrGroup
		groupRoutes.GET("getHmdrGroupList", groupController.List)                     // 获取HmdrGroup列表
	}

	hostRoutes := rg.Group("hmdrHost").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler()).Use(middleware.OperationRecord())
	{
		hostController := host.NewController(svcIns)

		hostRoutes.POST("createHmdrHost", hostController.Create)                  // 新建HmdrHost
		hostRoutes.DELETE("deleteHmdrHost", hostController.Delete)                // 删除HmdrHost
		hostRoutes.DELETE("deleteHmdrHostByIds", hostController.DeleteCollection) // 批量删除HmdrHost
		hostRoutes.PUT("updateHmdrHost", hostController.Update)                   // 更新HmdrHost
		hostRoutes.GET("findHmdrHost", hostController.Get)                        // 根据ID获取HmdrHost
		hostRoutes.GET("getHmdrHostList", hostController.List)                    // 获取HmdrHost列表
	}

	webSrvConfRoutes := rg.Group("conf").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler()).Use(middleware.OperationRecord())
	{
		webSrvConfController := web_server_config.NewController(svcIns)

		webSrvConfRoutes.GET("getOptions", webSrvConfController.GetOptions)                 // 获取options选择参数信息
		webSrvConfRoutes.POST("getConfInfo", webSrvConfController.GetConfigTextLines)       // 获取配置文件信息
		webSrvConfRoutes.POST("get-context-text", webSrvConfController.GetContextTextLines) // 获取上下文配置明细
		webSrvConfRoutes.POST("get-conf-struct", webSrvConfController.GetConfig)            // 获取配置文件JSON数据
		webSrvConfRoutes.POST("get-includes", webSrvConfController.GetIncludedConfigs)      // 获取包含的配置文件路径列表
		webSrvConfRoutes.POST("insert-clone-ctx", webSrvConfController.InsertWithClone)
		webSrvConfRoutes.POST("insert-new-ctx", webSrvConfController.InsertWithNew)
		webSrvConfRoutes.DELETE("remove-ctx", webSrvConfController.Remove)
		webSrvConfRoutes.POST("modify-ctx-value", webSrvConfController.ModifyContextValue)
		webSrvConfRoutes.POST("modify-clone-ctx", webSrvConfController.ModifyWithClone)
		webSrvConfRoutes.POST("change-ctx-enabled-state", webSrvConfController.ChangeContextEnabledState)
		webSrvConfRoutes.POST("modify-new-ctx", webSrvConfController.ModifyWithNew)
		webSrvConfRoutes.POST("move-ctx", webSrvConfController.Move)
	}

	webSrvBinCMDRoutes := rg.Group("bin-cmd").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler()).Use(middleware.OperationRecord())
	{
		webSrvBinCMDController := web_server_bin_cmd.NewController(svcIns)

		webSrvBinCMDRoutes.POST("exec", webSrvBinCMDController.Exec) // 提交Web服务端二进制命令工具执行请求
	}

	webSrvStatisticsRoutes := rg.Group("hmdr-statistics").Use(middleware.OperationRecord())
	{
		webSrvStatisticsController := web_server_statistics.NewController(svcIns)

		webSrvStatisticsRoutes.POST("proxy-svc-brief", webSrvStatisticsController.GetProxyServiceInfo) // 获取代理配置信息
	}
}
