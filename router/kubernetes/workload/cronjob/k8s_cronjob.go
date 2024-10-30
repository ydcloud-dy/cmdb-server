package cronjob

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sCronJobRouter struct{}

func (s *K8sCronJobRouter) Initk8sCronJobRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sCronJobRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sCronJobRouterWithoutRecord := Router.Group("kubernetes")
	//k8sCronJobRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sCronJobApi = v1.ApiGroupApp.CronJob.K8sCronJobApi
	{
		k8sCronJobRouter.POST("cronJob", k8sCronJobApi.CreateCronJob)   // 新建k8sCluster表
		k8sCronJobRouter.DELETE("cronJob", k8sCronJobApi.DeleteCronJob) // 删除k8sCluster表
		//k8sCronJobRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds) // 批量删除k8sCluster表
		k8sCronJobRouter.PUT("cronJob", k8sCronJobApi.UpdateCronJob)
	}
	{
		k8sCronJobRouterWithoutRecord.GET("cronJob", k8sCronJobApi.GetCronJobList) // 获取CronJob列表
		k8sCronJobRouterWithoutRecord.GET("cronJobDetails", k8sCronJobApi.DescribeCronJobInfo)
	}
	{
	}
}
