package roles

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sRoleRouter struct{}

func (s *K8sRoleRouter) Initk8sRoleRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sRoleRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sRoleRouterWithoutRecord := Router.Group("kubernetes")
	//k8sRoleRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sRoleApi = v1.ApiGroupApp.Role.K8sRolesApi
	{
		k8sRoleRouter.POST("Role", k8sRoleApi.CreateRoles)   // 新建k8sCluster表
		k8sRoleRouter.DELETE("Role", k8sRoleApi.DeleteRoles) // 删除k8sCluster表
		//k8sRoleRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds) // 批量删除k8sCluster表
		k8sRoleRouter.PUT("Role", k8sRoleApi.UpdateRoles)
	}
	{
		k8sRoleRouterWithoutRecord.GET("Role", k8sRoleApi.GetRolesList) // 获取Role列表
		k8sRoleRouterWithoutRecord.GET("RoleDetails", k8sRoleApi.DescribeRolesInfo)
	}
	{
	}
}
