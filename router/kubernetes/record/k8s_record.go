package record

import (
	v1 "DYCLOUD/api/v1"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sRecordRouter struct{}

func (s *K8sRecordRouter) Initk8sRecordRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sRecordRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sRecordRouterWithoutRecord := Router.Group("kubernetes")
	//k8sRecordRouterWithoutAuth := PublicRouter.Group("kubernetes")
	var k8sRecordApi = v1.ApiGroupApp.Record.K8sRecordApi
	{
		k8sRecordRouter.POST("Record", k8sRecordApi.CreateRecord)   // 新建k8sCluster表
		k8sRecordRouter.DELETE("Record", k8sRecordApi.DeleteRecord) // 删除k8sCluster表
		//k8sRecordRouter.DELETE("clusterByIds", k8sClusterApi.DeleteK8sClusterByIds) // 批量删除k8sCluster表
		k8sRecordRouter.PUT("Record", k8sRecordApi.UpdateRecord)
	}
	{
		k8sRecordRouterWithoutRecord.GET("Record", k8sRecordApi.GetRecordList) // 获取Record列表
		k8sRecordRouterWithoutRecord.GET("RecordDetails", k8sRecordApi.DescribeRecordInfo)
	}
	{
	}
}
