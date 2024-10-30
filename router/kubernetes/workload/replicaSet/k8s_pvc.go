package replicaSet

import (
	v1 "DYCLOUD/api/v1"
	"github.com/gin-gonic/gin"
)

type K8sReplicaSetRouter struct{}

func (s *K8sReplicaSetRouter) Initk8sReplicaSetRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	//k8sNodeRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sReplicaSetRouterWithoutRecord := Router.Group("kubernetes")
	//k8sNodeRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sReplicaSetApi = v1.ApiGroupApp.Replicaset.K8sReplicaSetApi
	{
		//k8sNodeRouter.POST("cluster", k8sClusterApi.CreateK8sCluster)             // 新建k8sCluster表
		//k8sNodeRouter.DELETE("cluster", k8sClusterApi.DeleteK8sCluster)           // 删除k8sCluster表
		//k8sNodeRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds) // 批量删除k8sCluster表

		//k8sNodeRouter.POST("clusters/detailById", k8sClusterApi.DescribeClusterById)
	}
	{
		//k8sNodeRouterWithoutRecord.POST("clusterById", k8sClusterApi.FindK8sCluster)   // 根据ID获取k8sCluster表
		k8sReplicaSetRouterWithoutRecord.GET("replicaSet", k8sReplicaSetApi.GetReplicaSetList) // 获取node列表
		//k8sReplicaSetRouterWithoutRecord.GET("ReplicaSets/metrics", k8sReplicaSetApi.MetricsReplicaSetsList)
		//k8sReplicaSetRouterWithoutRecord.GET("ReplicaSetDetails", k8sReplicaSetApi.DescribeReplicaSetInfo)
		//k8sNodeRouterWithoutRecord.GET("nodes/metrics", k8sReplicaSetApi.GetNodeMetricsList)
	}
	{
		//k8sNodeRouterWithoutAuth.GET("clusterPublic", k8sClusterApi.GetK8sClusterPublic) // 获取k8sCluster表列表
	}
}
