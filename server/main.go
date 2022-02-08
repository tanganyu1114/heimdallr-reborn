package main

import (
	"database/sql"
	"gin-vue-admin/core"
	"gin-vue-admin/global"
	"gin-vue-admin/initialize"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
	"gin-vue-admin/internal/hmdr_api/store/v1/bifrosts"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// @title Swagger Example API
// @version 0.0.1
// @description This is a sample Server pets
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-token
// @BasePath /
func main() {
	global.GVA_VP = core.Viper()          // 初始化Viper
	global.GVA_LOG = core.Zap()           // 初始化zap日志库
	global.GVA_DB = initialize.Gorm()     // gorm连接数据库
	initialize.MysqlTables(global.GVA_DB) // 初始化表
	// 程序结束前关闭数据库链接
	db, _ := global.GVA_DB.DB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			global.GVA_LOG.Error("failed to close database connection", zap.String("err", err.Error()))
		}
	}(db)

	// 初始化bifrost客户端仓库
	bifrostsStore := bifrosts.GetBifrostsStore()
	// 结束进程前关闭bifrost客户端仓库
	defer func(storeIns storev1.Factory) {
		err := storeIns.Close()
		if err != nil {
			global.GVA_LOG.Error("failed to close bifrost store", zap.String("err", err.Error()))
		}
	}(bifrostsStore)
	// 定时任务
	crontab := cron.New()
	// cron 定时同步客户端状态信息
	_, err := crontab.AddFunc("* * * * *", bifrostsStore.AgentInfos().SyncAgentInfos)
	if err != nil {
		global.GVA_LOG.Error("add crontab func err", zap.String("err", err.Error()))
	}
	crontab.Start()
	defer crontab.Stop()
	core.RunWindowsServer()
}
