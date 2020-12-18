package service

import (
	"context"
	"encoding/json"
	"gin-vue-admin/global"
	"go.uber.org/zap"
)

type Info struct {
	GroupName string     `json:"name"`
	Hosts     []HostInfo `json:"hosts"`
}

type HostInfo struct {
	Name      string    `json:"name"`
	Ipaddr    string    `json:"ipaddr"`
	Descrip   string    `json:"descrip"`
	Status    bool      `json:"status"`
	AgentInfo AgentInfo `json:"agent"`
}

type AgentInfo struct {
	OS             string        `json:"system"`
	Time           string        `json:"time"`
	Cpu            string        `json:"cpu"`
	Mem            string        `json:"mem"`
	Disk           string        `json:"disk"`
	StatusList     []NginxStatus `json:"status_list"`
	BifrostVersion string        `json:"bifrost_version"`
}

type NginxStatus struct {
	Name    string `json:"name"`
	Status  int    `json:"status"`
	Version string `json:"version"`
}

var GlobalAgentInfo = make([]Info, 0)

// 获取dashboard页面展示信息数据
func GetAgentInfo() (*[]Info, error) {
	return &GlobalAgentInfo, nil
}

func CronAgentInfo() {
	result := make([]Info, 0)

	for _, group := range BifrostGroups {
		tmpGroup := Info{
			GroupName: (*group).HmdrGroup.Name,
			Hosts:     make([]HostInfo, 0),
		}
		for _, host := range (*group).Hosts {
			var data AgentInfo
			var status = true
			bt, err := host.Client.Status(context.Background(), host.HmdrHost.Token)
			if err != nil {
				global.GVA_LOG.Error("access the client status failed, hostIp:"+host.HmdrHost.Ipaddr, zap.String("err", err.Error()))
				status = false
			}
			_ = json.Unmarshal(bt, &data)
			tmpHost := HostInfo{
				Name:      host.HmdrHost.Name,
				Ipaddr:    host.HmdrHost.Ipaddr,
				Descrip:   host.HmdrHost.Description,
				Status:    status,
				AgentInfo: data,
			}
			tmpGroup.Hosts = append(tmpGroup.Hosts, tmpHost)
		}
		result = append(result, tmpGroup)
	}
	GlobalAgentInfo = result
}
