package router

import (
	api "DYCLOUD/api/v1/cloudCmdb"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type CloudPlatformRouter struct {
}

func (c *CloudPlatformRouter) InitCloudPlatformRouter(Router *gin.RouterGroup) {
	cloudPlatformRouter := Router.Group("cloud_platform").Use(middleware.OperationRecord())
	cloudPlatformRouterWithoutRecord := Router.Group("cloud_platform")
	cloudPlatformApi := api.ApiGroupApp.CloudPlatformApi
	{
		cloudPlatformRouter.POST("getById", cloudPlatformApi.GetCloudPlatformById)           // 单个数据获取
		cloudPlatformRouter.POST("create", cloudPlatformApi.CreateCloudPlatform)             // 创建
		cloudPlatformRouter.PUT("update", cloudPlatformApi.UpdateCloudPlatform)              // 更新
		cloudPlatformRouter.DELETE("delete", cloudPlatformApi.DeleteCloudPlatform)           // 删除
		cloudPlatformRouter.DELETE("deleteByIds", cloudPlatformApi.DeleteCloudPlatformByIds) // 批量删除
	}
	{
		cloudPlatformRouterWithoutRecord.POST("list", cloudPlatformApi.CloudPlatformList) // 分页获取列表
	}
}
