package main

import (
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"

	"DYCLOUD/core"
	"DYCLOUD/global"
	"DYCLOUD/initialize"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// @title                       Gin-Vue-Admin Swagger API接口文档
// @version                     v2.7.4
// @description                 使用gin+vue进行极速开发的全栈开发基础平台
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
// @BasePath                    /
func main() {
	global.DYCLOUD_VP = core.Viper() // 初始化Viper
	initialize.OtherInit()
	global.DYCLOUD_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.DYCLOUD_LOG)
	global.DYCLOUD_DB = initialize.Gorm() // gorm连接数据库
	initialize.Timer()
	initialize.DBList()
	if global.DYCLOUD_DB != nil {
		initialize.RegisterTables() // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.DYCLOUD_DB.DB()
		defer db.Close()
	}
	core.RunWindowsServer()
}
