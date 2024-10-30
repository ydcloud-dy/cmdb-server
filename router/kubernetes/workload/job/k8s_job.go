package job

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sJobRouter struct{}

func (s *K8sJobRouter) Initk8sJobRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sJobRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sJobRouterWithoutRecord := Router.Group("kubernetes")
	//k8sJobRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sJobApi = v1.ApiGroupApp.Job.K8sJobApi
	{
		k8sJobRouter.POST("job", k8sJobApi.CreateJob)   // 新建k8sCluster表
		k8sJobRouter.DELETE("job", k8sJobApi.DeleteJob) // 删除k8sCluster表
		//k8sJobRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds) // 批量删除k8sCluster表
		k8sJobRouter.PUT("job", k8sJobApi.UpdateJob)
	}
	{
		k8sJobRouterWithoutRecord.GET("job", k8sJobApi.GetJobList) // 获取Job列表
		k8sJobRouterWithoutRecord.GET("jobDetails", k8sJobApi.DescribeJobInfo)
	}
	{
	}
}
