package limitRange

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sLimitRangeRouter struct{}

func (s *K8sLimitRangeRouter) Initk8sLimitRangeRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sNodeRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sLimitRangeRouterWithoutRecord := Router.Group("kubernetes")
	//k8sNodeRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sLimitRangeApi = v1.ApiGroupApp.LimitRange.K8sLimitRangeApi
	{
		k8sNodeRouter.POST("limitRange", k8sLimitRangeApi.CreateLimitRange)
		k8sNodeRouter.DELETE("limitRange", k8sLimitRangeApi.DeleteLimitRange)
		//k8sNodeRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds)
		k8sNodeRouter.PUT("limitRange", k8sLimitRangeApi.UpdateLimitRange)
		//k8sNodeRouter.POST("clusters/detailById", k8sClusterApi.DescribeClusterById)
	}
	{
		//k8sNodeRouterWithoutRecord.POST("clusterById", k8sClusterApi.FindK8sCluster)
		k8sLimitRangeRouterWithoutRecord.GET("limitRange", k8sLimitRangeApi.GetLimitRangeList)
		k8sLimitRangeRouterWithoutRecord.GET("limitRangeDetails", k8sLimitRangeApi.DescribeLimitRangeInfo)
		//k8sNodeRouterWithoutRecord.GET("nodes/metrics", k8sLimitRangeApi.GetNodeMetricsList)
	}
	{
		//k8sNodeRouterWithoutAuth.GET("clusterPublic", k8sClusterApi.GetK8sClusterPublic) // 获取k8sCluster表列表
	}
}
