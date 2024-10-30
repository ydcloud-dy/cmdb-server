package system

import (
	"DYCLOUD/global"
)

type JwtBlacklist struct {
	global.DYCLOUD_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
