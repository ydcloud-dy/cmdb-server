package core

import (
	"DYCLOUD/global"
	"DYCLOUD/initialize"
	"DYCLOUD/service/system"
	"fmt"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.DYCLOUD_CONFIG.System.UseMultipoint || global.DYCLOUD_CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
		initialize.RedisList()
	}

	if global.DYCLOUD_CONFIG.System.UseMongo {
		err := initialize.Mongo.Initialization()
		if err != nil {
			zap.L().Error(fmt.Sprintf("%+v", err))
		}
	}
	// 从db加载jwt数据
	if global.DYCLOUD_DB != nil {
		system.LoadAll()
	}

	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.DYCLOUD_CONFIG.System.Addr)
	s := initServer(address, Router)

	global.DYCLOUD_LOG.Info("server run success on ", zap.String("address", address))

	global.DYCLOUD_LOG.Error(s.ListenAndServe().Error())
}
