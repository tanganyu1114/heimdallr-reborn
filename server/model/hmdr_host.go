// 自动生成模板HmdrHost
package model

import (
	"gin-vue-admin/global"
)

// 如果含有time.Time 请自行import time包
type HmdrHost struct {
	global.GVA_MODEL
	GroupId     uint   `json:"groupId" form:"groupId" gorm:"column:group_id;comment:;type:bigint;size:20;"`
	Name        string `json:"name" form:"name" gorm:"column:name;comment:;type:varchar(50);size:50;"`
	Description string `json:"description" form:"description" gorm:"column:description;comment:;type:varchar(255);size:255;"`
	Status      bool   `json:"status" form:"status" gorm:"column:status;comment:;type:tinyint;"`
	Ipaddr      string `json:"ipaddr" form:"ipaddr" gorm:"column:ipaddr;comment:;type:char;"`
	Port        string `json:"port" form:"port" gorm:"column:port;comment:;type:char;"`
	Token       string `json:"token" form:"token" gorm:"column:token;comment:;type:varchar(255);size:255;"`
	Sequence    uint   `json:"sequence" form:"sequence" gorm:"column:sequence;comment:;type:bigint(20);size:20;"`
}

func (HmdrHost) TableName() string {
	return "hmdr_host"
}
