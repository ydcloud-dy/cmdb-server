package node

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sNodeRouter struct{}

func (s *K8sNodeRouter) Initk8sNodeRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sNodeRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sNodeRouterWithoutRecord := Router.Group("kubernetes")
	//k8sNodeRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sNodeApi = v1.ApiGroupApp.NodeApiGroup.K8sNodeApi
	{
		//k8sNodeRouter.POST("cluster", k8sClusterApi.CreateK8sCluster)             // 新建k8sCluster表
		//k8sNodeRouter.DELETE("cluster", k8sClusterApi.DeleteK8sCluster)           // 删除k8sCluster表
		//k8sNodeRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds) // 批量删除k8sCluster表
		k8sNodeRouter.PUT("nodes", k8sNodeApi.UpdateNodeInfo)
		k8sNodeRouter.POST("nodes/EvictAllPod", k8sNodeApi.EvictAllNodePod)
		//k8sNodeRouter.POST("clusters/detailById", k8sClusterApi.DescribeClusterById)
	}
	{
		//k8sNodeRouterWithoutRecord.POST("clusterById", k8sClusterApi.FindK8sCluster)   // 根据ID获取k8sCluster表
		k8sNodeRouterWithoutRecord.GET("nodes", k8sNodeApi.GetNodeList) // 获取node列表
		k8sNodeRouterWithoutRecord.GET("nodes/metrics", k8sNodeApi.GetNodeMetricsList)
		k8sNodeRouterWithoutRecord.GET("nodeDetails", k8sNodeApi.DescribeNodeInfo)
	}
	{
		//k8sNodeRouterWithoutAuth.GET("clusterPublic", k8sClusterApi.GetK8sClusterPublic) // 获取k8sCluster表列表
	}
}
