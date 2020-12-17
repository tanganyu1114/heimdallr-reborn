package service

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/model/request"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateHmdrHost
//@description: 创建HmdrHost记录
//@param: hmdrHost model.HmdrHost
//@return: err error

func CreateHmdrHost(hmdrHost model.HmdrHost) (err error) {
	err = global.GVA_DB.Create(&hmdrHost).Error
	if err == nil {
		CreateBifrostHost(hmdrHost)
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteHmdrHost
//@description: 删除HmdrHost记录
//@param: hmdrHost model.HmdrHost
//@return: err error

func DeleteHmdrHost(hmdrHost model.HmdrHost) (err error) {
	err = global.GVA_DB.Delete(hmdrHost).Error
	if err == nil {
		DeleteBifrostHost(hmdrHost)
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteHmdrHostByIds
//@description: 批量删除HmdrHost记录
//@param: ids request.IdsReq
//@return: err error

func DeleteHmdrHostByIds(ids request.IdsReq) (err error) {
	var hmdrHost []model.HmdrHost
	err = global.GVA_DB.Find(&hmdrHost, "id in ?", ids.Ids).Error
	err = global.GVA_DB.Delete(&[]model.HmdrHost{}, "id in ?", ids.Ids).Error
	if err == nil {
		for _, host := range hmdrHost {
			DeleteBifrostHost(host)
		}
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateHmdrHost
//@description: 更新HmdrHost记录
//@param: hmdrHost *model.HmdrHost
//@return: err error

func UpdateHmdrHost(hmdrHost *model.HmdrHost) (err error) {
	var oldHmdrHost model.HmdrHost
	global.GVA_DB.Find(&oldHmdrHost, hmdrHost.ID)
	err = global.GVA_DB.Save(hmdrHost).Error
	if err == nil {
		UpdateBifrostHost(oldHmdrHost, *hmdrHost)
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetHmdrHost
//@description: 根据id获取HmdrHost记录
//@param: id uint
//@return: err error, hmdrHost model.HmdrHost

func GetHmdrHost(id uint) (err error, hmdrHost model.HmdrHost) {
	err = global.GVA_DB.Where("id = ?", id).First(&hmdrHost).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetHmdrHostInfoList
//@description: 分页获取HmdrHost记录
//@param: info request.HmdrHostSearch
//@return: err error, list interface{}, total int64

func GetHmdrHostInfoList(info request.HmdrHostSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&model.HmdrHost{})
	var hmdrHosts []model.HmdrHost
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.GroupId != 0 {
		db = db.Where("`group_id` = ?", info.GroupId)
	}
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Order("sequence").Find(&hmdrHosts).Error
	return err, hmdrHosts, total
}
