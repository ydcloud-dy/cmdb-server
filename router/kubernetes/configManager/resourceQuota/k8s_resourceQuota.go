package resourceQuota

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sResourceQuotaRouter struct{}

func (s *K8sResourceQuotaRouter) Initk8sResourceQuotaRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sResourceQuotaRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sResourceQuotaRouterWithoutRecord := Router.Group("kubernetes")
	//k8sNodeRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sResourceQuotaApi = v1.ApiGroupApp.ResourceQuota.K8sResourceQuotaApi
	{
		k8sResourceQuotaRouter.POST("ResourceQuotas", k8sResourceQuotaApi.CreateResourceQuota)
		k8sResourceQuotaRouter.DELETE("ResourceQuotas", k8sResourceQuotaApi.DeleteResourceQuota)
		k8sResourceQuotaRouter.PUT("ResourceQuotas", k8sResourceQuotaApi.UpdateResourceQuota)
	}
	{

		k8sResourceQuotaRouterWithoutRecord.GET("ResourceQuotas", k8sResourceQuotaApi.GetResourceQuotaList) // 获取node列表
		k8sResourceQuotaRouterWithoutRecord.GET("ResourceQuotaDetails", k8sResourceQuotaApi.DescribeResourceQuotaInfo)

	}
	{

	}
}
