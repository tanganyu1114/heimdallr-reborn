package model

import "github.com/ClessLi/bifrost/pkg/client/bifrost"

type SocketControl struct {
	SearchConf
	LogName string `json:"log_name"`
	Status  bool   `json:"status"`
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
	LogWatcher bifrost.LogWatcher
	//DataChannel   chan []byte
	//SignalChannel chan int
	//ErrorChannel  chan error
}
