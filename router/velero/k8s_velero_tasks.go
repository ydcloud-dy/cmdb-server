package velero

import (
	"DYCLOUD/api/v1/velero"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type K8sVeleroTasksRouter struct{}

// InitK8sVeleroTasksRouter 初始化 k8sVeleroTasks表 路由信息
func (s *K8sVeleroTasksRouter) InitK8sVeleroTasksRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sVeleroTasksRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sVeleroTasksRouterWithoutRecord := Router.Group("kubernetes")
	//k8sVeleroTasksRouterWithoutAuth := PublicRouter.Group("kubernetes")
	api := velero.K8sVeleroTasksApi{}
	{
		k8sVeleroTasksRouter.POST("velero/tasks", api.CreateK8sVeleroTasks)
		k8sVeleroTasksRouter.POST("velero/record", api.CreateK8sVeleroRecord)
		k8sVeleroTasksRouter.DELETE("velero/tasks", api.DeleteK8sVeleroTasks)
		k8sVeleroTasksRouter.PUT("velero/tasks", api.UpdateK8sVeleroTasks)
		k8sVeleroTasksRouter.POST("velero", api.CreateVelero)
		k8sVeleroTasksRouter.DELETE("velero/record", api.DeleteK8sVeleroRecord)
		k8sVeleroTasksRouter.POST("velero/record/reduction", api.ReductionK8sVeleroRecord)
		k8sVeleroTasksRouter.DELETE("velero/restore", api.DeleteK8sVeleroRestore)

	}
	{
		k8sVeleroTasksRouterWithoutRecord.GET("velero/record", api.GetK8sVeleroRecordList)
		k8sVeleroTasksRouterWithoutRecord.GET("velero/recordDetail", api.DescribeK8sVeleroRecord)
		k8sVeleroTasksRouterWithoutRecord.GET("velero/restore", api.GetK8sVeleroRestoreList)
		k8sVeleroTasksRouterWithoutRecord.GET("velero/restoreDetail", api.DescribeK8sVeleroRestore)
		k8sVeleroTasksRouterWithoutRecord.GET("velero/taskDetail", api.DescribeVeleroTasks)
		k8sVeleroTasksRouterWithoutRecord.GET("velero/tasks", api.GetK8sVeleroTasksList)
	}
	{
		//k8sVeleroTasksRouterWithoutAuth.GET("getK8sVeleroTasksPublic", api.GetK8sVeleroTasksPublic) // 获取k8sVeleroTasks表列表
	}
}
