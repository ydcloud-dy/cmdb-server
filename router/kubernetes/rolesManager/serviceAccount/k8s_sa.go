package serviceAccount

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sServiceAccountRouter struct{}

func (s *K8sServiceAccountRouter) Initk8sServiceAccountRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sNodeRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sServiceAccountRouterWithoutRecord := Router.Group("kubernetes")
	//k8sNodeRouterWithoutAuth := PublicRouter.Group("kubernetes")

	var k8sServiceAccountApi = v1.ApiGroupApp.ServiceAccount.K8sServiceAccountApi
	{
		k8sNodeRouter.POST("serviceAccount", k8sServiceAccountApi.CreateServiceAccount)
		k8sNodeRouter.DELETE("serviceAccount", k8sServiceAccountApi.DeleteServiceAccount)
		k8sNodeRouter.PUT("serviceAccount", k8sServiceAccountApi.UpdateServiceAccount)
		//k8sNodeRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds) // 批量删除k8sCluster表

		//k8sNodeRouter.POST("clusters/detailById", k8sClusterApi.DescribeClusterById)
	}
	{
		//k8sNodeRouterWithoutRecord.POST("clusterById", k8sClusterApi.FindK8sCluster)   // 根据ID获取k8sCluster表
		k8sServiceAccountRouterWithoutRecord.GET("serviceAccount", k8sServiceAccountApi.GetServiceAccount)                 // 获取node列表
		k8sServiceAccountRouterWithoutRecord.GET("serviceAccountDetails", k8sServiceAccountApi.DescribeServiceAccountInfo) // 获取node列表

		//k8sServiceAccountRouterWithoutRecord.GET("ServiceAccounts/metrics", k8sServiceAccountApi.MetricsServiceAccountsList)
		//k8sServiceAccountRouterWithoutRecord.GET("ServiceAccountDetails", k8sServiceAccountApi.DescribeServiceAccountInfo)
		//k8sNodeRouterWithoutRecord.GET("nodes/metrics", k8sServiceAccountApi.GetNodeMetricsList)
	}
	{
		//k8sNodeRouterWithoutAuth.GET("clusterPublic", k8sClusterApi.GetK8sClusterPublic) // 获取k8sCluster表列表
	}
}
