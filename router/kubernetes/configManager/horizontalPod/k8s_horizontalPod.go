package horizontalPod

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sHorizontalPodRouter struct{}

func (s *K8sHorizontalPodRouter) Initk8sHorizontalPodRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sHorizontalPodRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sHorizontalPodRouterWithoutRecord := Router.Group("kubernetes")
	//k8sHorizontalPodRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sHorizontalPodApi = v1.ApiGroupApp.HorizontalPod.K8sHorizontalApi
	{
		k8sHorizontalPodRouter.POST("horizontalPod", k8sHorizontalPodApi.CreateHorizontal)   // 新建k8sCluster表
		k8sHorizontalPodRouter.DELETE("horizontalPod", k8sHorizontalPodApi.DeleteHorizontal) // 删除k8sCluster表
		//k8sHorizontalPodRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds) // 批量删除k8sCluster表
		k8sHorizontalPodRouter.PUT("horizontalPod", k8sHorizontalPodApi.UpdateHorizontal)
	}
	{
		k8sHorizontalPodRouterWithoutRecord.GET("horizontalPod", k8sHorizontalPodApi.GetHorizontalList) // 获取HorizontalPod列表
		k8sHorizontalPodRouterWithoutRecord.GET("horizontalPodDetail", k8sHorizontalPodApi.DescribeHorizontalInfo)
	}
	{
	}
}
