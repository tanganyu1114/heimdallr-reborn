package v1

import (
	"gin-vue-admin/global"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type Host struct {
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

type HostList struct {
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*Host `json:"items"`
}

func (h Host) TableName() string {
	return "hmdr_host"
}

func (h Host) GetOrder() uint {
	return h.Sequence
}

func (h Host) Key() interface{} {
	return h.ID
}
