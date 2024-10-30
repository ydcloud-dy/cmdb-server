package router

import (
	api "DYCLOUD/api/v1/cloudCmdb"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type CloudRDSRouter struct {
}

func (r *CloudRDSRouter) InitRDSRouter(Router *gin.RouterGroup) {
	cloudrdsRouter := Router.Group("rds").Use(middleware.OperationRecord())
	cloudrdsRouterWithoutRecord := Router.Group("rds")
	cloudrdsApi := api.ApiGroupApp.CloudRDSApi
	{
		cloudrdsRouter.POST("sync", cloudrdsApi.CloudRDSSync) // 同步RDS
		cloudrdsRouter.POST("tree", cloudrdsApi.CloudRDSTree) // 目录树
	}

	{
		cloudrdsRouterWithoutRecord.POST("list", cloudrdsApi.CloudRDSList) // 分页获取列表
	}
}
