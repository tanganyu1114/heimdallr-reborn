package v1

import (
	"gin-vue-admin/global"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type Group struct {
	global.GVA_MODEL
	Name        string `json:"name" form:"name" gorm:"column:name;comment:;type:varchar(50);size:50;"`
	Description string `json:"description" form:"description" gorm:"column:description;comment:;type:varchar(255);size:255;"`
	Sequence    uint   `json:"sequence" form:"sequence" gorm:"column:sequence;comment:;type:bigint(20);size:20;"`
}

type GroupList struct {
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*Group `json:"items"`
}

func (g Group) TableName() string {
	return "hmdr_group"
}

func (g Group) GetOrder() uint {
	return g.Sequence
}

func (g Group) Key() interface{} {
	return g.ID
}
