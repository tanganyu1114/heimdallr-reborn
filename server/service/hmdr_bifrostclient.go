package service

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/pkg/sort_map"
	"github.com/ClessLi/bifrost/pkg/client/bifrost"
	"go.uber.org/zap"
)

//var BifrostGroups = make(BifrostGroups,0)

//var BifrostGroups = make(map[uint]*model.BifrostGroup)
var BifrostGroups = model.NewBifrostGroups()

func SelectBifrostGroup(sf model.SearchConf) (*model.BifrostGroup, *model.BifrostHost) {
	group := BifrostGroups.Get(sf.GroupId)
	if group == nil {
		return nil, nil
	}
	host := group.(sort_map.SortMap).Get(sf.HostId)
	if host == nil {
		return group.(*model.BifrostGroup), nil
	}
	return group.(*model.BifrostGroup), host.(*model.BifrostHost)
}

// 初始化bifrost客户端
func InitBifrostClient() {

	// 查询分组信息
	db_group := global.GVA_DB.Model(&model.HmdrGroup{})
	var hmdrGroups []model.HmdrGroup
	res_group := db_group.Order("sequence").Find(&hmdrGroups)
	if res_group.Error != nil {
		global.GVA_LOG.Error("select hmdr_group error", zap.String("err", res_group.Error.Error()))
	}
	for _, group := range hmdrGroups {
		bGroup := model.NewBifrostGroup(group)
		BifrostGroups.Insert(&group, bGroup)
		// 查询主机信息h
		db := global.GVA_DB.Model(&model.HmdrHost{})
		var hmdrHost []model.HmdrHost
		db.Where("group_id = ? and status = 1", group.ID).Order("sequence").Find(&hmdrHost)

		// 初始化客户端
		for _, host := range hmdrHost {
			bifrostClient, initErr := bifrost.NewClient(host.Ipaddr + ":" + host.Port)
			if initErr != nil {
				global.GVA_LOG.Error("init bifrostClient Failed", zap.String("err", initErr.Error()))
				continue
			}

			bGroup.Hosts.Insert(&host, &model.BifrostHost{
				HmdrHost: host,
				Client:   bifrostClient,
			})
		}
	}
	// 调用查询客户端信息
	CronAgentInfo()
}

func CreateBifrostGroup(group model.HmdrGroup) {
	BifrostGroups.Insert(&group, model.NewBifrostGroup(group))
}

func CreateBifrostHost(host model.HmdrHost) {
	bifrostClient, initErr := bifrost.NewClient(host.Ipaddr + ":" + host.Port)
	if initErr != nil {
		global.GVA_LOG.Error("init bifrostClient Failed", zap.String("err", initErr.Error()))
	}
	bGroup := BifrostGroups.Get(host.GroupId)
	if bGroup == nil {
		return
	}
	bGroup.(*model.BifrostGroup).Hosts.Insert(&host, &model.BifrostHost{
		HmdrHost: host,
		Client:   bifrostClient,
	})
}

func DeleteBifrostHost(host model.HmdrHost) {
	bGroup := BifrostGroups.Get(host.GroupId)
	if bGroup != nil {
		bHost := bGroup.(*model.BifrostGroup).Hosts.Get(host.ID)
		if bHost != nil {
			// 关闭客户端
			bHost.(*model.BifrostHost).Client.Close()
			bGroup.(*model.BifrostGroup).Hosts.Remove(host.ID)
		}
	}
}

func DeleteBifrostGroup(group model.HmdrGroup) {
	bGroup := BifrostGroups.Get(group.ID)

	if bGroup != nil {
		closeGourp := func(k, v interface{}) bool {
			key := k.(uint)
			value := v.(*model.BifrostHost)
			value.Client.Close()
			bGroup.(*model.BifrostGroup).Hosts.Remove(key)
			return true
		}
		// 删除组之前 关闭清理所有组下的客户端
		//for u, host := range bgroup.Hosts {
		//	host.Client.Close()
		//	delete(bgroup.Hosts, u)
		//}
		bGroup.(*model.BifrostGroup).Hosts.Range(closeGourp)
		//delete(BifrostGroups, group.ID)
		BifrostGroups.Remove(group.ID)
	}
}

func UpdateBifrostHost(oldhost, newhost model.HmdrHost) {
	DeleteBifrostHost(oldhost)
	if newhost.Status {
		CreateBifrostHost(newhost)
	}
}
