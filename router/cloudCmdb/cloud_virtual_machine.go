package router

import (
	api "DYCLOUD/api/v1/cloudCmdb"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type CloudVirtualMachineRouter struct {
}

func (c *CloudVirtualMachineRouter) InitVirtualMachineRouter(Router *gin.RouterGroup) {
	cloudvirtualMachineRouter := Router.Group("virtualMachine").Use(middleware.OperationRecord())
	cloudvirtualMachineRouterWithoutRecord := Router.Group("virtualMachine")
	cloudvirtualMachineApi := api.ApiGroupApp.CloudVirtualMachineApi
	{
		cloudvirtualMachineRouter.POST("sync", cloudvirtualMachineApi.CloudVirtualMachineSync) // 同步云主机
		cloudvirtualMachineRouter.POST("tree", cloudvirtualMachineApi.CloudVirtualMachineTree) // 目录树
	}

	{
		cloudvirtualMachineRouterWithoutRecord.POST("list", cloudvirtualMachineApi.CloudVirtualMachineList) // 分页获取列表
	}
}
