package router

import (
	api "DYCLOUD/api/v1/cloudCmdb"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type CloudRegionRouter struct {
}

func (c *CloudRegionRouter) InitCloudRegionRouter(Router *gin.RouterGroup) {
	cloudRegionRouter := Router.Group("cloud_region").Use(middleware.OperationRecord())
	cloudRegionApi := api.ApiGroupApp.CloudRegionApi
	{
		cloudRegionRouter.POST("syncRegion", cloudRegionApi.CloudPlatformSyncRegion) // 同步区域信息
	}
}
