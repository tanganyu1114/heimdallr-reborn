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
	"time"

	"github.com/gin-gonic/gin"
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

		groupRoutes.
			POST("createHmdrGroup", groupController.Create).                   // 新建HmdrGroup
			DELETE("deleteHmdrGroup", groupController.Delete).                 // 删除HmdrGroup
			DELETE("deleteHmdrGroupByIds", groupController.DeleteCollections). // 批量删除HmdrGroup
			PUT("updateHmdrGroup", groupController.Update).                    // 更新HmdrGroup
			GET("findHmdrGroup", groupController.Get).                         // 根据ID获取HmdrGroup
			GET("getHmdrGroupList", groupController.List)                      // 获取HmdrGroup列表
	}

	hostRoutes := rg.Group("hmdrHost").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler()).Use(middleware.OperationRecord())
	{
		hostController := host.NewController(svcIns)

		hostRoutes.
			POST("createHmdrHost", hostController.Create).                  // 新建HmdrHost
			DELETE("deleteHmdrHost", hostController.Delete).                // 删除HmdrHost
			DELETE("deleteHmdrHostByIds", hostController.DeleteCollection). // 批量删除HmdrHost
			PUT("updateHmdrHost", hostController.Update).                   // 更新HmdrHost
			GET("findHmdrHost", hostController.Get).                        // 根据ID获取HmdrHost
			GET("getHmdrHostList", hostController.List)                     // 获取HmdrHost列表
	}

	webSrvConfRoutes := rg.Group("conf").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler()).Use(middleware.OperationRecord())
	{
		webSrvConfController := web_server_config.NewController(svcIns)

		webSrvConfRoutes.
			GET("getOptions", webSrvConfController.GetOptions).                    // 获取options选择参数信息
			POST("getConfInfo", webSrvConfController.GetConfigTextLines).          // 获取配置文件信息
			POST("get-context-text", webSrvConfController.GetContextTextLines).    // 获取上下文配置明细
			POST("get-conf-struct", webSrvConfController.GetConfig).               // 获取配置文件JSON数据
			POST("get-includes", webSrvConfController.GetIncludedConfigs).         // 获取包含的配置文件路径列表
			POST("search-ctx-poses", webSrvConfController.SearchContextPositions). // 搜索上下文坐标列表
			POST("insert-clone-ctx", webSrvConfController.InsertWithClone).
			POST("insert-new-ctx", webSrvConfController.InsertWithNew).
			DELETE("remove-ctx", webSrvConfController.Remove).
			POST("update-conf", webSrvConfController.UpdateConfig).
			POST("modify-ctx-value", webSrvConfController.ModifyContextValue).
			POST("modify-clone-ctx", webSrvConfController.ModifyWithClone).
			POST("change-ctx-enabled-state", webSrvConfController.ChangeContextEnabledState).
			POST("modify-new-ctx", webSrvConfController.ModifyWithNew).
			POST("move-ctx", webSrvConfController.Move)
	}

	webSrvBinCMDRoutes := rg.Group("bin-cmd").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler()).Use(middleware.OperationRecord())
	{
		webSrvBinCMDController := web_server_bin_cmd.NewController(svcIns)

		webSrvBinCMDRoutes.POST("exec", webSrvBinCMDController.Exec) // 提交Web服务端二进制命令工具执行请求
	}

	webSrvStatisticsRoutes := rg.Group("hmdr-statistics").Use(middleware.OperationRecord())
	{
		webSrvStatisticsController := web_server_statistics.NewController(svcIns)

		webSrvStatisticsRoutes.POST("proxy-svc-brief", webSrvStatisticsController.GetProxyServiceInfo)                     // 获取代理配置信息
		webSrvStatisticsRoutes.POST("conn-check-of-proxy-svc", webSrvStatisticsController.ConnectivityCheckOfProxyService) // 代理服务网络连通性检查
		webSrvStatisticsRoutes.POST("export-proxy-svc-excel", webSrvStatisticsController.ExportProxyServiceInfoToExcel)    // 导出代理信息为Excel
	}
}
