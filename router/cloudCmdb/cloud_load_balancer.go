package router

import (
	api "DYCLOUD/api/v1/cloudCmdb"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type CloudLoadBalancerRouter struct {
}

func (l *CloudLoadBalancerRouter) InitLoadBalancerRouter(Router *gin.RouterGroup) {
	cloudloadBalancerRouter := Router.Group("loadBalancer").Use(middleware.OperationRecord())
	cloudloadBalancerRouterWithoutRecord := Router.Group("loadBalancer")
	cloudloadBalancerApi := api.ApiGroupApp.CloudLoadBalancerApi
	{
		cloudloadBalancerRouter.POST("sync", cloudloadBalancerApi.CloudLoadBalancerSync) // 同步负载均衡
		cloudloadBalancerRouter.POST("tree", cloudloadBalancerApi.CloudLoadBalancerTree) // 目录树
	}

	{
		cloudloadBalancerRouterWithoutRecord.POST("list", cloudloadBalancerApi.CloudLoadBalancerList) // 分页获取列表
	}
}
