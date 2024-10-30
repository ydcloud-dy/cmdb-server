package configmap

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sConfigMapRouter struct{}

// InitK8sClusterRouter 初始化 k8sCluster表 路由信息
func (s *K8sConfigMapRouter) InitK8sConfigMapRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sConfigMapRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sConfigMapRouterWithoutRecord := Router.Group("kubernetes")
	//k8sConfigMapRouterWithoutAuth := PublicRouter.Group("kubernetes")

	var k8sConfigMapApi = v1.ApiGroupApp.ConfigMapGroup.K8sConfigMapApi
	{
		k8sConfigMapRouter.POST("configMap", k8sConfigMapApi.CreateConfigMap)
		k8sConfigMapRouter.DELETE("configMap", k8sConfigMapApi.DeleteConfigMap)

		k8sConfigMapRouter.PUT("configMap", k8sConfigMapApi.UpdateConfigMap)

	}
	{
		k8sConfigMapRouterWithoutRecord.GET("configMap", k8sConfigMapApi.GetConfigMapList)             // 根据ID获取k8sCluster表
		k8sConfigMapRouterWithoutRecord.GET("configMapDetails", k8sConfigMapApi.DescribeConfigMapInfo) // 根据ID获取k8sCluster表

		//k8sConfigMapRouterWithoutRecord.GET("ConfigMap/detail", k8sConfigMapApi.DescribeConfigMapInfo)

	}
	{

	}
}
