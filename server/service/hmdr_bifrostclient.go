package service

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"github.com/ClessLi/bifrost/pkg/client/bifrost"
	"go.uber.org/zap"
)

//var BifrostGroups = make(BifrostGroups,0)

var BifrostGroups = make(map[uint]*model.BifrostGroup)

func SelectBifrostGroup(sf model.SearchConf) bool {
	_, ok := BifrostGroups[sf.GroupId]
	if !ok {
		return ok
	}
	_, ok = BifrostGroups[sf.GroupId].Hosts[sf.HostId]
	if !ok {
		return ok
	}
	return ok
}

// 初始化bifrost客户端
func InitBifrostClient() {

	// 查询分组信息
	db_group := global.GVA_DB.Model(&model.HmdrGroup{})
	var hmdrGroups []model.HmdrGroup
	res_group := db_group.Find(&hmdrGroups)
	if res_group.Error != nil {
		global.GVA_LOG.Error("select hmdr_group error", zap.String("err", res_group.Error.Error()))
	}
	for _, group := range hmdrGroups {
		BifrostGroups[group.ID] = &model.BifrostGroup{
			HmdrGroup: group,
			//ID: group.ID,
			Hosts: make(map[uint]*model.BifrostHost),
		}

		// 查询主机信息h
		db := global.GVA_DB.Model(&model.HmdrHost{})
		var hmdrHost []model.HmdrHost
		db.Where("group_id = ? and status = 1", group.ID).Find(&hmdrHost)

		// 初始化客户端
		for _, host := range hmdrHost {
			bifrostClient, initErr := bifrost.NewClient(host.Ipaddr + ":" + host.Port)
			if initErr != nil {
				global.GVA_LOG.Error("init bifrostClient Failed", zap.String("err", initErr.Error()))
				continue
			}
			BifrostGroups[group.ID].Hosts[host.ID] = &model.BifrostHost{
				HmdrHost: host,
				Client:   bifrostClient,
			}
		}
	}
}

func CreateBifrostGroup(group model.HmdrGroup) {
	BifrostGroups[group.ID] = &model.BifrostGroup{
		HmdrGroup: group,
		Hosts:     make(map[uint]*model.BifrostHost),
	}
}

func CreateBifrostHost(host model.HmdrHost) {
	bifrostClient, initErr := bifrost.NewClient(host.Ipaddr + ":" + host.Port)
	if initErr != nil {
		global.GVA_LOG.Error("init bifrostClient Failed", zap.String("err", initErr.Error()))
	}
	BifrostGroups[host.GroupId].Hosts[host.ID] = &model.BifrostHost{
		HmdrHost: host,
		Client:   bifrostClient,
	}
}

func DeleteBifrostHost(host model.HmdrHost) {
	group, ok := BifrostGroups[host.GroupId]
	if ok {
		ghost, ok := group.Hosts[host.ID]
		if ok {
			// 关闭客户端
			ghost.Client.Close()
			delete(group.Hosts, host.ID)
		}
	}
}

func DeleteBifrostGroup(group model.HmdrGroup) {
	bgroup, ok := BifrostGroups[group.ID]
	if ok {
		// 删除组之前 关闭清理所有组下的客户端
		for u, host := range bgroup.Hosts {
			host.Client.Close()
			delete(bgroup.Hosts, u)
		}
		delete(BifrostGroups, group.ID)
	}
}

func UpdateBifrostHost(oldhost, newhost model.HmdrHost) {
	DeleteBifrostHost(oldhost)
	if newhost.Status {
		CreateBifrostHost(newhost)
	}
}
