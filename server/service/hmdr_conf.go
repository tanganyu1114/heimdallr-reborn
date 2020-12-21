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

func GetOptions() ([]Group, error) {
	reslut := make([]Group, 0)

	getOptFromGroup := func(k, v interface{}) bool {
		group := v.(*model.BifrostGroup)
		tmpGroup := Group{
			Label:    group.HmdrGroup.Name,
			Value:    group.HmdrGroup.ID,
			Children: make([]Host, 0),
		}
		getOptFromHost := func(k, v interface{}) bool {
			host := v.(*model.BifrostHost)
			tmpHost := Host{
				Label:    host.HmdrHost.Name,
				Value:    host.HmdrHost.ID,
				Children: make([]Srv, 0),
			}
			var data AgentInfo
			bt, err := host.Client.Status(context.Background(), host.HmdrHost.Token)
			if err != nil {
				global.GVA_LOG.Error("Get SrvName Failed", zap.String("err", err.Error()))
				//return false  // 注意：如果报错，散列表后续元素将不被加载
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
			return true
		}
		group.Hosts.Range(getOptFromHost)
		reslut = append(reslut, tmpGroup)
		return true
	}
	BifrostGroups.Range(getOptFromGroup)
	return reslut, nil
}

func GetConfInfo(sf model.SearchConf) *[]byte {

	_, bHost := SelectBifrostGroup(sf)
	if bHost != nil {
		data, err := bHost.Client.ViewConfig(context.Background(), bHost.HmdrHost.Token, sf.SrvName)
		if err != nil {
			global.GVA_LOG.Error("Get Conf Failed", zap.String("err", err.Error()))
			return nil
		}
		return &data
	}
	return nil
}
