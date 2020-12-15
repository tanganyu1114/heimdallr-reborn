package service

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

var GlobalSocketInfo = make(map[uint]*model.SocketGroup)

func SelectLogWatcher(sc model.SocketControl) bool {
	_, ok := GlobalSocketInfo[sc.GroupId]
	if !ok {
		return ok
	}
	_, ok = GlobalSocketInfo[sc.GroupId].Host[sc.HostId]
	if !ok {
		return ok
	}
	_, ok = GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName]
	if !ok {
		return ok
	}
	_, ok = GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName]
	if !ok {
		return ok
	}
	return ok
}

func CreateLogWatcher(sc model.SocketControl) {
	// 判断客户端是否存在
	if ok := SelectBifrostGroup(sc.SearchConf); ok {
		logWatcher, err := BifrostGroups[sc.GroupId].Hosts[sc.HostId].Client.WatchLog(context.Background(), BifrostGroups[sc.GroupId].Hosts[sc.HostId].HmdrHost.Token, sc.SrvName, sc.LogName)
		if err != nil {
			global.GVA_LOG.Error("init logwatcher failed", zap.Any("err", err))
		}
		GlobalSocketInfo[sc.GroupId] = &model.SocketGroup{
			Host: make(map[uint]*model.SocketHost),
		}
		GlobalSocketInfo[sc.GroupId].Host[sc.HostId] = &model.SocketHost{
			SrvName: make(map[string]*model.SocketSrvName),
		}
		GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName] = &model.SocketSrvName{
			LogName: make(map[string]*model.SocketInfo),
		}
		GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName] = &model.SocketInfo{
			LogWatcher: logWatcher,
			Count:      1,
		}
	}
}

func CalcLogWatcherCount(sc model.SocketControl, count int) {

	GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].Count += count
	if GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].Count == 0 {
		DeleteLogWater(sc)
	}
}

func DeleteLogWater(sc model.SocketControl) {
	if ok := SelectLogWatcher(sc); ok {
		GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].LogWatcher.Close()
		delete(GlobalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName, sc.LogName)
	}
}
