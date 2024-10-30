package deployment

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sDeploymentRouter struct{}

// InitK8sClusterRouter 初始化 k8sCluster表 路由信息
func (s *K8sDeploymentRouter) InitK8sDeploymentRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sDeploymentRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sDeploymentRouterWithoutRecord := Router.Group("kubernetes")
	//k8sDeploymentRouterWithoutAuth := PublicRouter.Group("kubernetes")

	var k8sDeploymentApi = v1.ApiGroupApp.DeploymentGroup.K8sDeploymentApi
	{

		k8sDeploymentRouter.PUT("deployment", k8sDeploymentApi.UpdateDeploymentInfo)
		k8sDeploymentRouter.POST("deployment", k8sDeploymentApi.CreateDeployment)
		k8sDeploymentRouter.DELETE("deployment", k8sDeploymentApi.DeleteDeployment)
		k8sDeploymentRouter.PATCH("deployment", k8sDeploymentApi.RollBackDeployment)

	}
	{
		k8sDeploymentRouterWithoutRecord.GET("deployment", k8sDeploymentApi.GetDeploymentList) // 根据ID获取k8sCluster表
		k8sDeploymentRouterWithoutRecord.GET("deployment/detail", k8sDeploymentApi.DescribeDeploymentInfo)

	}
	{

	}
}
