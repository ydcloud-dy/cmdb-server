package Ingress

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sIngressRouter struct{}

func (s *K8sIngressRouter) Initk8sIngressRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sNodeRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sIngressRouterWithoutRecord := Router.Group("kubernetes")
	//k8sNodeRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sIngressApi = v1.ApiGroupApp.Ingress.K8sIngressApi
	{
		k8sNodeRouter.POST("ingress", k8sIngressApi.CreateIngress)
		k8sNodeRouter.DELETE("ingress", k8sIngressApi.DeleteIngress)
		//k8sNodeRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds)
		k8sNodeRouter.PUT("ingress", k8sIngressApi.UpdateIngress)
		//k8sNodeRouter.POST("clusters/detailById", k8sClusterApi.DescribeClusterById)
	}
	{
		//k8sNodeRouterWithoutRecord.POST("clusterById", k8sClusterApi.FindK8sCluster)
		k8sIngressRouterWithoutRecord.GET("ingress", k8sIngressApi.GetIngressList)
		k8sIngressRouterWithoutRecord.GET("ingressDetails", k8sIngressApi.DescribeIngressInfo)
		//k8sNodeRouterWithoutRecord.GET("nodes/metrics", k8sIngressApi.GetNodeMetricsList)
	}
	{
		//k8sNodeRouterWithoutAuth.GET("clusterPublic", k8sClusterApi.GetK8sClusterPublic) // 获取k8sCluster表列表
	}
}
