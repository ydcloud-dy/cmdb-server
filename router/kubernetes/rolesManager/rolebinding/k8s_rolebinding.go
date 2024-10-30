package rolebinding

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sRoleBindingRouter struct{}

func (s *K8sRoleBindingRouter) Initk8sRoleBindingRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sRoleBindingRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sRoleBindingRouterWithoutRecord := Router.Group("kubernetes")
	//k8sRoleBindingRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sRoleBindingApi = v1.ApiGroupApp.RoleBinding.K8sRoleBindingApi
	{
		k8sRoleBindingRouter.POST("RoleBinding", k8sRoleBindingApi.CreateRoleBinding)   // 新建k8sCluster表
		k8sRoleBindingRouter.DELETE("RoleBinding", k8sRoleBindingApi.DeleteRoleBinding) // 删除k8sCluster表
		//k8sRoleBindingRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds) // 批量删除k8sCluster表
		k8sRoleBindingRouter.PUT("RoleBinding", k8sRoleBindingApi.UpdateRoleBinding)
	}
	{
		k8sRoleBindingRouterWithoutRecord.GET("RoleBinding", k8sRoleBindingApi.GetRoleBindingList) // 获取RoleBinding列表
		k8sRoleBindingRouterWithoutRecord.GET("RoleBindingDetails", k8sRoleBindingApi.DescribeRoleBindingInfo)
	}
	{
	}
}
