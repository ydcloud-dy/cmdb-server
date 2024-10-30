package cluster

import (
	"DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sClusterRouter struct{}

// InitK8sClusterRouter 初始化 k8sCluster表 路由信息
func (s *K8sClusterRouter) InitK8sClusterRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sClusterRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sClusterRouterWithoutRecord := Router.Group("kubernetes")
	k8sClusterRouterWithoutAuth := PublicRouter.Group("kubernetes")

	var k8sClusterApi = v1.ApiGroupApp.ClusterApiGroup.K8sClusterApi
	{
		k8sClusterRouter.POST("cluster", k8sClusterApi.CreateK8sCluster)             // 新建k8sCluster表
		k8sClusterRouter.DELETE("cluster", k8sClusterApi.DeleteK8sCluster)           // 删除k8sCluster表
		k8sClusterRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds) // 批量删除k8sCluster表
		k8sClusterRouter.PUT("cluster", k8sClusterApi.UpdateK8sCluster)              // 更新k8sCluster表
		k8sClusterRouter.POST("credential", k8sClusterApi.CreateCredential)          // 创建凭据
		k8sClusterRouter.POST("getUserById", k8sClusterApi.GetClusterUserById)       // 集群用户获取
		k8sClusterRouter.POST("getClusterRoles", k8sClusterApi.GetClusterRoles)
		k8sClusterRouter.POST("getClusterApiGroups", k8sClusterApi.GetClusterApiGroups)         // 资源分组获取
		k8sClusterRouter.POST("createClusterRole", k8sClusterApi.CreateClusterRole)             // 创建角色
		k8sClusterRouter.PUT("updateClusterRole", k8sClusterApi.UpdateClusterRole)              // 更新角色
		k8sClusterRouter.DELETE("deleteClusterRole", k8sClusterApi.DeleteClusterRole)           // 删除角色
		k8sClusterRouter.POST("createClusterUser", k8sClusterApi.CreateClusterUser)             //  创建用户授权
		k8sClusterRouter.PUT("updateClusterUser", k8sClusterApi.UpdateClusterUser)              //  更新用户授权
		k8sClusterRouter.DELETE("deleteClusterUser", k8sClusterApi.DeleteClusterUser)           //  删除用户授权
		k8sClusterRouter.POST("getClusterUserNamespace", k8sClusterApi.GetClusterUserNamespace) //  用户获取命名空间
		k8sClusterRouter.POST("getClusterListNamespace", k8sClusterApi.GetClusterListNamespace) //  管理员获取命名空间
		//k8sClusterRouter.POST("getUserNamespaceList", k8sClusterApi.GetClusterUserNamespaceListForNS)
		//k8sClusterRouter.GET("cloudtty", k8sClusterApi.)
		//k8sClusterRouter.POST("clusters/detailById", k8sClusterApi.DescribeClusterById)
	}
	{
		k8sClusterRouterWithoutRecord.POST("clusterById", k8sClusterApi.FindK8sCluster)   // 根据ID获取k8sCluster表
		k8sClusterRouterWithoutRecord.GET("clusterList", k8sClusterApi.GetK8sClusterList) // 获取k8sCluster表列表
	}
	{
		k8sClusterRouterWithoutAuth.GET("clusterPublic", k8sClusterApi.GetK8sClusterPublic) // 获取k8sCluster表列表
	}
}
