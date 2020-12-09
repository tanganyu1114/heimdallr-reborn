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

// 获取dashboard页面展示信息数据
func GetAgentInfo() (*[]Info, error) {

	/*	result := make([]Info, 0)
		// 查询分组信息
		// 创建db
		db_group := global.GVA_DB.Model(&model.HmdrGroup{})
		var hmdrGroups []model.HmdrGroup
		res_group := db_group.Find(&hmdrGroups)
		if res_group.Error != nil {
			global.GVA_LOG.Error("select hmdr_group error", zap.String("err", res_group.Error.Error()))
			return nil, res_group.Error
		}
		// 查询对应组别下的主机信息
		// 创建db
		// 存在分组信息才查询主机
		if res_group.RowsAffected > 0 {
			for _, group := range hmdrGroups {
				db_host := global.GVA_DB.Model(&model.HmdrHost{})
				var hmdrHosts []model.HmdrHost
				info := new(Info)
				// 设置组名
				info.GroupName = group.Name
				db_host.Where("group_id = ?", group.ID).Find(&hmdrHosts)
				for _, host := range hmdrHosts {
					// 如果host是启用的，则获取agent服务器相关信息
					var data AgentInfo
					if host.Status {
						// data = GetBifrost(&host)
						bc := SelectBifrostClient(model.GroupAndHost{GroupId: host.GroupId,HostId: host.ID})
						bt,err := bc.Client.Status(context.Background(),bc.Token)
						if err != nil {
							fmt.Println(err)
						}
						json.Unmarshal(bt, &data)
					}
					info.Hosts = append(info.Hosts, HostInfo{
						Name:      host.Name,
						Status:    host.Status,
						Ipaddr:    host.Ipaddr,
						Descrip:   host.Description,
						AgentInfo: data,
					})
				}
				result = append(result, *info)
			}
		} else {
			return nil, nil
		}
		return &result, nil*/

	// 全局客户端
	/*	result := make([]Info, 0)
		result = append(result, Info{GroupName: })
		for _, c := range initialize.BifrostClients {

			var data AgentInfo

			bt,err:=c.Client.Status(context.Background(),c.Token)

			if err != nil {
				global.GVA_LOG.Error("get client status faild", zap.String("err", err.Error()))
				continue
			}
			err = json.Unmarshal(bt, &data)
			if err != nil {
				global.GVA_LOG.Error("json unmarshal faild !", zap.String("err", err.Error()))
				continue
			}
			info := new(Info)
			info.Hosts = append(info.Hosts, HostInfo{
				Name:      host.Name,
				Status:    host.Status,
				Ipaddr:    host.Ipaddr,
				Descrip:   host.Description,
				AgentInfo: data,
			})
			return
		}*/
	result := make([]Info, 0)
	for _, group := range BifrostGroups {
		tmpGroup := Info{
			GroupName: (*group).HmdrGroup.Name,
			Hosts:     make([]HostInfo, 0),
		}
		for _, host := range (*group).Hosts {
			var data AgentInfo
			bt, err := host.Client.Status(context.Background(), host.HmdrHost.Token)
			if err != nil {
				global.GVA_LOG.Error("access the client status failed", zap.String("err", err.Error()))
			}
			_ = json.Unmarshal(bt, &data)
			tmpHost := HostInfo{
				Name:      host.HmdrHost.Name,
				Ipaddr:    host.HmdrHost.Ipaddr,
				Descrip:   host.HmdrHost.Description,
				Status:    true,
				AgentInfo: data,
			}
			tmpGroup.Hosts = append(tmpGroup.Hosts, tmpHost)
		}
		result = append(result, tmpGroup)
	}
	return &result, nil
}

/*func GetBifrost(host *model.HmdrHost) (data AgentInfo) {
	// data = new(AgentInfo)
	bifrostClient, initErr := bifrost.NewClient(host.Ipaddr + ":" + host.Port)
	if initErr != nil {
		global.GVA_LOG.Error("init bifrostClient Failed", zap.String("err", initErr.Error()))
		return
	}
	defer bifrostClient.Close()
	jsondata, err := bifrostClient.Status(context.Background(), host.Token)
	if err != nil {
		global.GVA_LOG.Error("init bifrostClient Failed", zap.String("err", err.Error()))
		return
	}
	err = json.Unmarshal(jsondata, &data)
	if err != nil {
		global.GVA_LOG.Error("json unmarshal faild !", zap.String("err", err.Error()))
		return
	}
	return data
}*/
