package service

import (
	"context"
	"encoding/json"
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"go.uber.org/zap"
)

type Group struct {
	Label    string `json:"label"`
	Value    uint   `json:"value"`
	Children []Host `json:"children"`
}
type Host struct {
	Label    string `json:"label"`
	Value    uint   `json:"value"`
	Children []Srv  `json:"children"`
}
type Srv struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

/*func GetOptions()([]Group,error)  {
	reslut := make([]Group,0)
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
		db_host := global.GVA_DB.Model(&model.HmdrHost{})
		var hmdrHosts []model.HmdrHost

		for _, group := range hmdrGroups {
			tmpGroup := Group{
				Label:    group.Name,
				Value:    group.ID,
				Children: make([]Host, 0),
			}
			db_host.Where("group_id = ?", group.ID).Find(&hmdrHosts)
			for _, host := range hmdrHosts {
				tmpHost := Host{
					Label:    host.Name,
					Value:    host.ID,
					Children: make([]Srv, 0),
				}
				var data AgentInfo
				if host.Status {
					data = GetBifrost(&host)
				}
				for _, val := range data.StatusList {
					tmpSrv := Srv{
						Label: val.Name,
						Value: val.Name,
					}
					tmpHost.Children = append(tmpHost.Children, tmpSrv)
				}
				tmpGroup.Children = append(tmpGroup.Children, tmpHost)
			}
			reslut = append(reslut, tmpGroup)
		}
		return  reslut,nil
	}else {
		return nil,nil
	}
}*/

func GetOptions() ([]Group, error) {
	reslut := make([]Group, 0)

	for _, group := range BifrostGroups {
		tmpGroup := Group{
			Label:    group.HmdrGroup.Name,
			Value:    group.HmdrGroup.ID,
			Children: make([]Host, 0),
		}
		for _, host := range (*group).Hosts {
			tmpHost := Host{
				Label:    host.HmdrHost.Name,
				Value:    host.HmdrHost.ID,
				Children: make([]Srv, 0),
			}
			var data AgentInfo
			bt, err := host.Client.Status(context.Background(), host.HmdrHost.Token)
			if err != nil {
				global.GVA_LOG.Error("Get SrvName Failed", zap.String("err", err.Error()))
			}
			_ = json.Unmarshal(bt, &data)
			for _, val := range data.StatusList {
				tmpSrv := Srv{
					Label: val.Name,
					Value: val.Name,
				}
				tmpHost.Children = append(tmpHost.Children, tmpSrv)
			}
			tmpGroup.Children = append(tmpGroup.Children, tmpHost)
		}
		reslut = append(reslut, tmpGroup)
	}
	return reslut, nil
}

func GetConfInfo(sf model.SearchConf) *[]byte {

	/*	// 查询主机信息
		db_host := global.GVA_DB.Model(&model.HmdrHost{})
		var hmdrHost model.HmdrHost
		db_host.Where("id = ? AND group_id = ?",sf.HostId,sf.GroupId).First(&hmdrHost)

		// 获取配置文件信息
		bifrostClient, initErr := bifrost.NewClient(hmdrHost.Ipaddr + ":" + hmdrHost.Port)
		if initErr != nil {
			global.GVA_LOG.Error("init bifrostClient Failed", zap.String("err", initErr.Error()))
			return nil,initErr
		}
		defer bifrostClient.Close()
		data,err := bifrostClient.ViewConfig(context.Background(),hmdrHost.Token,sf.SrvName)

		return &data,err*/

	group, ok := BifrostGroups[sf.GroupId]
	if ok {
		host, ok := group.Hosts[sf.HostId]
		if ok {
			data, err := host.Client.ViewConfig(context.Background(), host.HmdrHost.Token, sf.SrvName)
			if err != nil {
				global.GVA_LOG.Error("Get Conf Failed", zap.String("err", err.Error()))
			}
			return &data
		}
		return nil
	}
	return nil
}
