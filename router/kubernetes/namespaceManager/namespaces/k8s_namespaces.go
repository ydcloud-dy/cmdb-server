package namespaces

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sNamespaceRouter struct{}

func (s *K8sNamespaceRouter) Initk8sNamespaceRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sNamespaceRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sNamespaceRouterWithoutRecord := Router.Group("kubernetes")
	//k8sNamespaceRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sNamespaceApi = v1.ApiGroupApp.Namespace.K8sNamespaceApi
	{
		k8sNamespaceRouter.POST("namespace", k8sNamespaceApi.CreateNamespace)   // 新建k8sCluster表
		k8sNamespaceRouter.DELETE("namespace", k8sNamespaceApi.DeleteNamespace) // 删除k8sCluster表
		//k8sNamespaceRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds) // 批量删除k8sCluster表
		k8sNamespaceRouter.PUT("namespace", k8sNamespaceApi.UpdateNamespace)
	}
	{
		k8sNamespaceRouterWithoutRecord.GET("namespace", k8sNamespaceApi.GetNamespaceList) // 获取Namespace列表
		k8sNamespaceRouterWithoutRecord.GET("namespaceDetails", k8sNamespaceApi.DescribeNamespaceInfo)
	}
	{
	}
}
