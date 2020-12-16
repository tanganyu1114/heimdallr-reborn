package main

import (
	"gin-vue-admin/core"
	"gin-vue-admin/global"
	"gin-vue-admin/initialize"
	"gin-vue-admin/service"
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
	defer db.Close()
	// 初始化bifrost客户端
	service.InitBifrostClient()
	// 关闭所有的bifrost客户端
	defer func() {
		for _, bg := range service.BifrostGroups {
			for _, c := range (*bg).Hosts {
				c.Client.Close()
			}
		}
	}()
	// 定时任务
	crontab := cron.New()
	// cron 获取客户端信息
	_, err := crontab.AddFunc("* * * * *", service.CronAgentInfo)
	if err != nil {
		global.GVA_LOG.Error("add crontab func err", zap.String("err", err.Error()))
	}
	crontab.Start()
	defer crontab.Stop()
	core.RunWindowsServer()
}
