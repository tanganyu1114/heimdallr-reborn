package datas

import (
	"os"
	"time"

	"github.com/gookit/color"

	"github.com/tanganyu1114/heimdallr-reborn/server/model"

	"gorm.io/gorm"
)

var Authorities = []model.SysAuthority{
	{CreatedAt: time.Now(), UpdatedAt: time.Now(), AuthorityId: "888", AuthorityName: "普通用户", ParentId: "0"},
	{CreatedAt: time.Now(), UpdatedAt: time.Now(), AuthorityId: "8881", AuthorityName: "普通用户子角色", ParentId: "888"},
	{CreatedAt: time.Now(), UpdatedAt: time.Now(), AuthorityId: "9528", AuthorityName: "测试角色", ParentId: "0"},
	{CreatedAt: time.Now(), UpdatedAt: time.Now(), AuthorityId: "8882", AuthorityName: "SDK用户", ParentId: "0"},
}

func InitSysAuthority(db *gorm.DB) {
	if err := db.Transaction(func(tx *gorm.DB) error {
		for _, auth := range Authorities {
			var existing model.SysAuthority
			result := tx.Where("authority_id = ?", auth.AuthorityId).First(&existing)
			if result.Error != nil {
				if err := tx.Create(&auth).Error; err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		color.Warn.Printf("[Mysql]--> sys_authorities 表的初始数据失败,err: %v\n", err)
		os.Exit(0)
	}
	color.Info.Println("[Mysql]--> sys_authorities 表的初始数据成功")
}
