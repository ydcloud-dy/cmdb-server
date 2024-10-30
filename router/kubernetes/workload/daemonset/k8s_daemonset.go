package daemonset

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sDaemonSetRouter struct{}

func (s *K8sDaemonSetRouter) Initk8sDaemonSetRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sDaemonSetRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sDaemonSetRouterWithoutRecord := Router.Group("kubernetes")
	//k8sDaemonSetRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sDaemonSetApi = v1.ApiGroupApp.DaemonSet.K8sDaemonSetApi
	{
		k8sDaemonSetRouter.PUT("daemonset", k8sDaemonSetApi.UpdateDaemonsSet)
		k8sDaemonSetRouter.DELETE("daemonset", k8sDaemonSetApi.DeleteDaemonSet)
		//k8sDaemonSetRouter.PUT("DaemonSets", k8sDaemonSetApi.UpdateDaemonSetInfo)
		k8sDaemonSetRouter.POST("daemonset", k8sDaemonSetApi.CreateDaemonSet)
		//k8sDaemonSetRouter.POST("clusters/detailById", k8sClusterApi.DescribeClusterById)
	}
	{
		//k8sDaemonSetRouterWithoutRecord.POST("clusterById", k8sClusterApi.FindK8sCluster)   // 根据ID获取k8sCluster表
		k8sDaemonSetRouterWithoutRecord.GET("daemonset", k8sDaemonSetApi.GetDaemonSetList) // 获取DaemonSet列表
		//k8sDaemonSetRouterWithoutRecord.GET("DaemonSets/metrics", k8sDaemonSetApi.GetDaemonSetMetricsList)
		k8sDaemonSetRouterWithoutRecord.GET("daemonsetDetails", k8sDaemonSetApi.DescribeDaemonSetInfo)
	}
	{
		//k8sDaemonSetRouterWithoutAuth.GET("clusterPublic", k8sClusterApi.GetK8sClusterPublic) // 获取k8sCluster表列表
	}
}
