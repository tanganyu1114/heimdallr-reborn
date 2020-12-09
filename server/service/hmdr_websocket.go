package service

import (
	"context"
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/utils"
	"github.com/ClessLi/bifrost/pkg/client/bifrost"
	"go.uber.org/zap"
)

var GolbalSocketInfo = make(map[string]model.GlobalSocketInfo)

func GetLogsInfo(sf model.SearchConf) string {
	// 查询主机信息h
	db_host := global.GVA_DB.Model(&model.HmdrHost{})
	var hmdrHost model.HmdrHost
	db_host.Where("id = ? AND group_id = ?", sf.HostId, sf.GroupId).First(&hmdrHost)

	// 初始化客户端
	bifrostClient, initErr := bifrost.NewClient(hmdrHost.Ipaddr + ":" + hmdrHost.Port)
	if initErr != nil {
		global.GVA_LOG.Error("init bifrostClient Failed", zap.String("err", initErr.Error()))
	}
	defer bifrostClient.Close()

	// 初始化全局socket信息
	key := utils.Struct2Md5(sf)
	GolbalSocketInfo[key] = model.GlobalSocketInfo{
		Count:         1,
		DataChannel:   make(chan []byte, 1),
		SignalChannel: make(chan int, 1),
	}

	// 获取日志信息
	//dataChan := make(chan []byte, 1)
	//signal := make(chan int, 1)
	//timeout := time.After(time.Second * 20)
	go func() {
		err := bifrostClient.WatchLog(context.Background(), hmdrHost.Token, sf.SrvName, "access.log", GolbalSocketInfo[key].DataChannel, GolbalSocketInfo[key].SignalChannel)
		if err != nil {
			global.GVA_LOG.Error("watch log err:", zap.String("err", err.Error()))
		}
	}()

	defer func() { GolbalSocketInfo[key].SignalChannel <- 9 }()

	return key
	//for {
	//	select {
	//	case data := <-dataChan:
	//
	//	case <-timeout:
	//
	//		return
	//	}
	//}

}
