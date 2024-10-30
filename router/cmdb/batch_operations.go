package cmdb

import (
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type BatchOperationsRouter struct{}

// InitBatchOperationsRouter 初始化 BatchOperations表 路由信息
func (s *BatchOperationsRouter) InitBatchOperationsRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	BatchOperationsRouter := Router.Group("cmdb").Use(middleware.OperationRecord())
	//BatchOperationsRouterWithoutRecord := Router.Group("cmdb")
	//BatchOperationsRouterWithoutAuth := PublicRouter.Group("cmdb")
	{
		BatchOperationsRouter.POST("batchOperations/execute", batchOperationsApi.ExecuteCommands) // 新建BatchOperations表
		BatchOperationsRouter.GET("batchOperations/execLogs",batchOperationsApi.ExecuteRecords)
	}
	{
		//BatchOperationsRouterWithoutRecord.GET("hostsById", batchOperationsApi.FindBatchOperations) // 根据ID获取BatchOperations表
		//BatchOperationsRouterWithoutRecord.GET("hosts", batchOperationsApi.GetBatchOperationsList)  // 获取BatchOperations表列表
	}
	{
		//BatchOperationsRouterWithoutAuth.GET("hostsPublic", batchOperationsApi.GetBatchOperationsPublic) // BatchOperations表开放接口
	}
}
