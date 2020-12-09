// 自动生成模板HmdrGroup
package model

import (
	"gin-vue-admin/global"
)

// 如果含有time.Time 请自行import time包
type HmdrGroup struct {
	global.GVA_MODEL
	Name        string `json:"name" form:"name" gorm:"column:name;comment:;type:varchar(50);size:50;"`
	Description string `json:"description" form:"description" gorm:"column:description;comment:;type:varchar(255);size:255;"`
}

func (HmdrGroup) TableName() string {
	return "hmdr_group"
}
