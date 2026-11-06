package router

import (
	"time"

	"github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/controller/v1/agent_info"
	"github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/controller/v1/group"
	"github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/controller/v1/host"
	"github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/controller/v1/web_server_bin_cmd"
	"github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/controller/v1/web_server_config"
	"github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/controller/v1/web_server_log_watcher"
	"github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/controller/v1/web_server_statistics"
	bifrostssvc "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/service/v1/bifrosts"
	loggingsvc "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/service/v1/logging"
	bifrostsstore "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1/bifrosts"
	cachestore "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1/cache"
	"github.com/tanganyu1114/heimdallr-reborn/server/middleware"

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

	webSrvConfRoutes := rg.Group("conf").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		webSrvConfController := web_server_config.NewController(svcIns)

		// 读接口：仅查询配置信息，无副作用，不记录操作日志
		webSrvConfRoutes.
			GET("getOptions", webSrvConfController.GetOptions).                   // 获取 options 选择参数信息
			POST("getConfInfo", webSrvConfController.GetConfigTextLines).         // 获取配置文件信息
			POST("get-context-text", webSrvConfController.GetContextTextLines).   // 获取上下文配置明细
			POST("get-conf-struct", webSrvConfController.GetConfig).              // 获取配置文件 JSON 数据
			POST("get-includes", webSrvConfController.GetIncludedConfigs).        // 获取包含的配置文件路径列表
			POST("search-ctx-poses", webSrvConfController.SearchContextPositions) // 搜索上下文坐标列表

		// 写接口：会变更配置文件内容，统一记录操作日志
		webSrvConfRoutes.Use(middleware.OperationRecord()).
			POST("insert-clone-ctx", webSrvConfController.InsertWithClone).                   // 克隆已有配置上下文后插入
			POST("insert-new-ctx", webSrvConfController.InsertWithNew).                       // 插入全新的配置上下文
			DELETE("remove-ctx", webSrvConfController.Remove).                                // 删除配置上下文
			POST("update-conf", webSrvConfController.UpdateConfig).                           // 更新配置文件
			POST("modify-ctx-value", webSrvConfController.ModifyContextValue).                // 修改配置上下文的值
			POST("modify-clone-ctx", webSrvConfController.ModifyWithClone).                   // 用克隆的配置上下文替换原有内容
			POST("change-ctx-enabled-state", webSrvConfController.ChangeContextEnabledState). // 启用或禁用配置上下文
			POST("modify-new-ctx", webSrvConfController.ModifyWithNew).                       // 用全新的配置上下文替换原有内容
			POST("move-ctx", webSrvConfController.Move)                                       // 调整配置上下文的位置
	}

	webSrvBinCMDRoutes := rg.Group("bin-cmd").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler()).Use(middleware.OperationRecord())
	{
		webSrvBinCMDController := web_server_bin_cmd.NewController(svcIns)

		webSrvBinCMDRoutes.POST("exec", webSrvBinCMDController.Exec) // 提交Web服务端二进制命令工具执行请求
	}

	webSrvStatisticsRoutes := rg.Group("hmdr-statistics").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		webSrvStatisticsController := web_server_statistics.NewController(svcIns)

		webSrvStatisticsRoutes.POST("proxy-svc-brief", webSrvStatisticsController.GetProxyServiceInfo)                     // 获取代理配置信息
		webSrvStatisticsRoutes.POST("conn-check-of-proxy-svc", webSrvStatisticsController.ConnectivityCheckOfProxyService) // 代理服务网络连通性检查
		webSrvStatisticsRoutes.POST("export-proxy-svc-excel", webSrvStatisticsController.ExportProxyServiceInfoToExcel)    // 导出代理信息为Excel
	}
}
