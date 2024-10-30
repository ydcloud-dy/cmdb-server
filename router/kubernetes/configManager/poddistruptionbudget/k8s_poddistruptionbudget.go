package poddistruptionbudget

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sPoddisruptionbudgetRouter struct{}

func (s *K8sPoddisruptionbudgetRouter) Initk8sPoddisruptionbudgetRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sPoddisruptionbudgetRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sPoddisruptionbudgetRouterWithoutRecord := Router.Group("kubernetes")
	//k8sPoddisruptionbudgetRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sPoddisruptionbudgetApi = v1.ApiGroupApp.Poddisruptionbudget.K8sPoddisruptionbudgetApi
	{
		k8sPoddisruptionbudgetRouter.POST("Poddisruptionbudget", k8sPoddisruptionbudgetApi.CreatePoddisruptionbudget)   // 新建k8sCluster表
		k8sPoddisruptionbudgetRouter.DELETE("Poddisruptionbudget", k8sPoddisruptionbudgetApi.DeletePoddisruptionbudget) // 删除k8sCluster表
		//k8sPoddisruptionbudgetRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds) // 批量删除k8sCluster表
		k8sPoddisruptionbudgetRouter.PUT("Poddisruptionbudget", k8sPoddisruptionbudgetApi.UpdatePoddisruptionbudget)
	}
	{
		k8sPoddisruptionbudgetRouterWithoutRecord.GET("Poddisruptionbudget", k8sPoddisruptionbudgetApi.GetPoddisruptionbudgetList) // 获取Poddisruptionbudget列表
		k8sPoddisruptionbudgetRouterWithoutRecord.GET("PoddisruptionbudgetDetails", k8sPoddisruptionbudgetApi.DescribePoddisruptionbudgetInfo)
	}
	{
	}
}
