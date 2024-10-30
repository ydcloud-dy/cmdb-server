package statefulSet

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sStatefulSetRouter struct{}

func (s *K8sStatefulSetRouter) Initk8sStatefulSetRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sStatefulSetRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sStatefulSetRouterWithoutRecord := Router.Group("kubernetes")
	//k8sStatefulSetRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sStatefulSetApi = v1.ApiGroupApp.StatefulSet.K8sStatefulSetApi
	{
		k8sStatefulSetRouter.POST("statefulset", k8sStatefulSetApi.CreateStatefulset)   // 新建k8sCluster表
		k8sStatefulSetRouter.DELETE("statefulset", k8sStatefulSetApi.DeleteStatefulSet) // 删除k8sCluster表
		//k8sStatefulSetRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds) // 批量删除k8sCluster表
		k8sStatefulSetRouter.PUT("statefulset", k8sStatefulSetApi.UpdateStatefulSetInfo)
		//k8sStatefulSetRouter.POST("StatefulSets/EvictAllPod", k8sStatefulSetApi.EvictAllStatefulSetPod)
		//k8sStatefulSetRouter.POST("clusters/detailById", k8sClusterApi.DescribeClusterById)
	}
	{
		//k8sStatefulSetRouterWithoutRecord.POST("clusterById", k8sClusterApi.FindK8sCluster)   // 根据ID获取k8sCluster表
		k8sStatefulSetRouterWithoutRecord.GET("statefulset", k8sStatefulSetApi.GetStatefulSetList) // 获取StatefulSet列表
		//k8sStatefulSetRouterWithoutRecord.GET("StatefulSets/metrics", k8sStatefulSetApi.GetStatefulSetMetricsList)
		k8sStatefulSetRouterWithoutRecord.GET("statefulsetDetails", k8sStatefulSetApi.DescribeStatefulSetInfo)
	}
	{
		//k8sStatefulSetRouterWithoutAuth.GET("clusterPublic", k8sClusterApi.GetK8sClusterPublic) // 获取k8sCluster表列表
	}
}
