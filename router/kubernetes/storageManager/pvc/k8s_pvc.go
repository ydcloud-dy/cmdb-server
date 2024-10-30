package pvc

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sPvcRouter struct{}

func (s *K8sPvcRouter) Initk8sPvcRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sNodeRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sPvcRouterWithoutRecord := Router.Group("kubernetes")
	//k8sNodeRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sPvcApi = v1.ApiGroupApp.PvcGroup.K8sPvcApi
	{
		k8sNodeRouter.POST("pvc", k8sPvcApi.CreatePVC)   // 新建k8sCluster表
		k8sNodeRouter.DELETE("pvc", k8sPvcApi.DeletePVC) // 删除k8sCluster表
		k8sNodeRouter.PUT("pvc", k8sPvcApi.UpdatePVC)    // 删除k8sCluster表

		//k8sNodeRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds) // 批量删除k8sCluster表

		//k8sNodeRouter.POST("clusters/detailById", k8sClusterApi.DescribeClusterById)
	}
	{
		//k8sNodeRouterWithoutRecord.POST("clusterById", k8sClusterApi.FindK8sCluster)   // 根据ID获取k8sCluster表
		k8sPvcRouterWithoutRecord.GET("pvc", k8sPvcApi.GetPvcList)             // 获取node列表
		k8sPvcRouterWithoutRecord.GET("pvcDetails", k8sPvcApi.DescribePVCInfo) // 获取node列表

		//k8sPvcRouterWithoutRecord.GET("Pvcs/metrics", k8sPvcApi.MetricsPvcsList)
		//k8sPvcRouterWithoutRecord.GET("PvcDetails", k8sPvcApi.DescribePvcInfo)
		//k8sNodeRouterWithoutRecord.GET("nodes/metrics", k8sPvcApi.GetNodeMetricsList)
	}
	{
		//k8sNodeRouterWithoutAuth.GET("clusterPublic", k8sClusterApi.GetK8sClusterPublic) // 获取k8sCluster表列表
	}
}
