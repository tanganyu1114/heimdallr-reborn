package model

import (
	"github.com/tanganyu1114/heimdallr-reborn/global"
)

type JwtBlacklist struct {
	global.GVA_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
