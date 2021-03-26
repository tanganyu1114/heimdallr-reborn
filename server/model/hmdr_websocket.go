package model

import (
	"gin-vue-admin/utils/pipe/log_watcher_pipe"
	"github.com/ClessLi/bifrost/pkg/client/bifrost"
)

type SocketControl struct {
	SearchConf
	LogName string `json:"log_name"`
}

type SocketGroup struct {
	Host map[uint]*SocketHost
}

type SocketHost struct {
	SrvName map[string]*SocketSrvName
}

type SocketSrvName struct {
	LogName map[string]*SocketInfo
}

type SocketInfo struct {
	Count      int
	LogWatcher bifrost.WatchClient
	LogPipe    log_watcher_pipe.LogWatcherPipe
}

func (sc SocketControl) GetLogWatcher(globalSocketInfo map[uint]*SocketGroup) bifrost.WatchClient {
	return globalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].LogWatcher
}

func (sc SocketControl) GetLogPipe(globalSocketInfo map[uint]*SocketGroup) log_watcher_pipe.LogWatcherPipe {
	return globalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].LogPipe
}

func (sc SocketControl) GetSocketCount(globalSocketInfo map[uint]*SocketGroup) int {
	return globalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].Count
}

func (sc SocketControl) InsertSocketCount(globalSocketInfo map[uint]*SocketGroup, count int) {
	globalSocketInfo[sc.GroupId].Host[sc.HostId].SrvName[sc.SrvName].LogName[sc.LogName].Count = count
}
