package datas

import (
	"github.com/tanganyu1114/heimdallr-reborn/global"
	"os"
	"time"

	"github.com/gookit/color"

	"github.com/tanganyu1114/heimdallr-reborn/model"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var Users = []model.SysUser{
	{GVA_MODEL: global.GVA_MODEL{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, UUID: uuid.NewV4(), Username: "admin", Password: "e10adc3949ba59abbe56e057f20f883e", NickName: "超级管理员", HeaderImg: "http://qmplusimg.henrongyi.top/gva_header.jpg", AuthorityId: "888"},
	{GVA_MODEL: global.GVA_MODEL{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, UUID: uuid.NewV4(), Username: "a303176530", Password: "3ec063004a6f31642261936a379fde3d", NickName: "QMPlusUser", HeaderImg: "http://qmplusimg.henrongyi.top/1572075907logo.png", AuthorityId: "9528"},
	{GVA_MODEL: global.GVA_MODEL{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}, UUID: uuid.NewV4(), Username: "sdk_user", Password: "e10adc3949ba59abbe56e057f20f883e", NickName: "SDK测试用户", HeaderImg: "http://qmplusimg.henrongyi.top/gva_header.jpg", AuthorityId: "8882"},
}

func InitSysUser(db *gorm.DB) {
	if err := db.Transaction(func(tx *gorm.DB) error {
		for _, user := range Users {
			var existing model.SysUser
			result := tx.Where("id = ?", user.ID).First(&existing)
			if result.Error != nil {
				if err := tx.Create(&user).Error; err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		color.Warn.Printf("[Mysql]--> sys_users 表的初始数据失败,err: %v\n", err)
		os.Exit(0)
	}
	color.Info.Println("[Mysql]--> sys_users 表的初始数据成功")
}
