package cloudtty

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sCloudTTYRouter struct{}

func (s *K8sCloudTTYRouter) Initk8sCloudTTYRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sCloudTTYRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	//k8sCloudTTYRouterWithoutRecord := Router.Group("kubernetes")
	//k8sCloudTTYRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sCloudTTYApi = v1.ApiGroupApp.CloudTTY.K8sCloudTTYApi
	{
		k8sCloudTTYRouter.POST("/cloudtty/get", k8sCloudTTYApi.CloudTTYGet) // CloudTTY
	}
	{
		//k8sCloudTTYRouterWithoutRecord.GET("cloudTTY", k8sCloudTTYApi.GetCloudTTYList) // 获取CloudTTY列表
		//k8sCloudTTYRouterWithoutRecord.GET("cloudTTYDetails", k8sCloudTTYApi.DescribeCloudTTYInfo)
	}
	{
	}
}
