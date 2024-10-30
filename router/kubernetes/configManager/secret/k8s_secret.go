package secret

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sSecretRouter struct{}

func (s *K8sSecretRouter) Initk8sSecretRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sNodeRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sSecretRouterWithoutRecord := Router.Group("kubernetes")
	//k8sNodeRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sSecretApi = v1.ApiGroupApp.Secret.K8sSecretApi
	{
		k8sNodeRouter.POST("secret", k8sSecretApi.CreateSecret)   // 新建k8sCluster表
		k8sNodeRouter.DELETE("secret", k8sSecretApi.DeleteSecret) // 新建k8sCluster表
		k8sNodeRouter.PUT("secret", k8sSecretApi.UpdateSecret)    // 新建k8sCluster表

		//k8sNodeRouter.DELETE("cluster", k8sClusterApi.DeleteK8sCluster)           // 删除k8sCluster表
		//k8sNodeRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds) // 批量删除k8sCluster表

		//k8sNodeRouter.POST("clusters/detailById", k8sClusterApi.DescribeClusterById)
	}
	{
		////k8sNodeRouterWithoutRecord.POST("clusterById", k8sClusterApi.FindK8sCluster)   // 根据ID获取k8sCluster表
		k8sSecretRouterWithoutRecord.GET("secret", k8sSecretApi.GetSecretList)             // 获取node列表
		k8sSecretRouterWithoutRecord.GET("secretDetails", k8sSecretApi.DescribeSecretInfo) // 获取node列表

		//k8sSecretRouterWithoutRecord.GET("Secrets/metrics", k8sSecretApi.MetricsSecretsList)
		//k8sSecretRouterWithoutRecord.GET("SecretDetails", k8sSecretApi.DescribeSecretInfo)
		//k8sNodeRouterWithoutRecord.GET("nodes/metrics", k8sSecretApi.GetNodeMetricsList)
	}
	{
		//k8sNodeRouterWithoutAuth.GET("clusterPublic", k8sClusterApi.GetK8sClusterPublic) // 获取k8sCluster表列表
	}
}
