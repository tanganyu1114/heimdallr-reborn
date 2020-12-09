package service

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/model/request"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateHmdrGroup
//@description: 创建HmdrGroup记录
//@param: hmdrGroup model.HmdrGroup
//@return: err error

func CreateHmdrGroup(hmdrGroup model.HmdrGroup) (err error) {
	err = global.GVA_DB.Create(&hmdrGroup).Error
	if err == nil {
		// 创建组
		CreateBifrostGroup(hmdrGroup)
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteHmdrGroup
//@description: 删除HmdrGroup记录
//@param: hmdrGroup model.HmdrGroup
//@return: err error

func DeleteHmdrGroup(hmdrGroup model.HmdrGroup) (err error) {
	err = global.GVA_DB.Delete(hmdrGroup).Error
	if err == nil {
		// 删除组
		DeleteBifrostGroup(hmdrGroup)
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteHmdrGroupByIds
//@description: 批量删除HmdrGroup记录
//@param: ids request.IdsReq
//@return: err error

func DeleteHmdrGroupByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]model.HmdrGroup{}, "id in ?", ids.Ids).Error
	if err == nil {
		for _, id := range ids.Ids {
			DeleteBifrostGroup(model.HmdrGroup{GVA_MODEL: global.GVA_MODEL{ID: uint(id)}})
		}

	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateHmdrGroup
//@description: 更新HmdrGroup记录
//@param: hmdrGroup *model.HmdrGroup
//@return: err error

func UpdateHmdrGroup(hmdrGroup *model.HmdrGroup) (err error) {
	err = global.GVA_DB.Save(hmdrGroup).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetHmdrGroup
//@description: 根据id获取HmdrGroup记录
//@param: id uint
//@return: err error, hmdrGroup model.HmdrGroup

func GetHmdrGroup(id uint) (err error, hmdrGroup model.HmdrGroup) {
	err = global.GVA_DB.Where("id = ?", id).First(&hmdrGroup).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetHmdrGroupInfoList
//@description: 分页获取HmdrGroup记录
//@param: info request.HmdrGroupSearch
//@return: err error, list interface{}, total int64

func GetHmdrGroupInfoList(info request.HmdrGroupSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&model.HmdrGroup{})
	var hmdrGroups []model.HmdrGroup
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Find(&hmdrGroups).Error
	return err, hmdrGroups, total
}
