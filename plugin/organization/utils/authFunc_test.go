package utils

import (
	"DYCLOUD/core"
	"DYCLOUD/global"
	"DYCLOUD/initialize"
	"testing"
)

func TestGetAllUserIDs(t *testing.T) {
	global.DYCLOUD_VP = core.Viper()      // 初始化Viper
	global.DYCLOUD_LOG = core.Zap()       // 初始化zap日志库
	global.DYCLOUD_DB = initialize.Gorm() // gorm连接数据库
}
