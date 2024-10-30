package pod

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sPodRouter struct{}

func (s *K8sPodRouter) Initk8sPodRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sPodRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sPodRouterWithoutRecord := Router.Group("kubernetes")
	//k8sNodeRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sPodApi = v1.ApiGroupApp.PodApiGroup.K8sPodApi
	{
		k8sPodRouter.POST("pods", k8sPodApi.CreatePod)
		k8sPodRouter.DELETE("pods", k8sPodApi.DeletePod)
		k8sPodRouter.PUT("pods", k8sPodApi.UpdatePod)
		k8sPodRouter.GET("events", k8sPodApi.PodEvents)
		k8sPodRouter.POST("pods/listFiles", k8sPodApi.ListPodFiles)
		k8sPodRouter.POST("pods/uploadFile", k8sPodApi.UploadFiles)
		k8sPodRouter.POST("pods/deleteFiles", k8sPodApi.DeleteFiles)

	}
	{

		k8sPodRouterWithoutRecord.GET("pods", k8sPodApi.GetPodList) // 获取node列表
		k8sPodRouterWithoutRecord.GET("pods/metrics", k8sPodApi.MetricsPodsList)
		k8sPodRouterWithoutRecord.GET("podDetails", k8sPodApi.DescribePodInfo)

	}
	{

	}
}
