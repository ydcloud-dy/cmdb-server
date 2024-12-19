package cicd

import (
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type PipelinesRouter struct{}

// InitPipelinesRouter 初始化 Pipelines表 路由信息
func (s *PipelinesRouter) InitPipelinesRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	PipelinesRouter := Router.Group("cicd").Use(middleware.OperationRecord())
	{
		PipelinesRouter.GET("pipelines", PipelineApi.GetPipelinesList)
		PipelinesRouter.POST("pipelines", PipelineApi.CreatePipelines)
		PipelinesRouter.PUT("pipelines", PipelineApi.UpdatePipelines)
		PipelinesRouter.DELETE("pipelines/:id", PipelineApi.DeletePipelines)
		PipelinesRouter.GET("pipelines/:id", PipelineApi.DescribePipelines)
		PipelinesRouter.DELETE("pipelines", PipelineApi.DeletePipelinesByIds)
		PipelinesRouter.GET("pipelinesStatus", PipelineApi.GetPipelinesStatus)
		PipelinesRouter.GET("pipelines/notice/:id", PipelineApi.GetPipelinesNotice)
		PipelinesRouter.POST("pipelines/closeNotice/:id", PipelineApi.ClosePipelineNotice)
		PipelinesRouter.POST("pipelines/notice", PipelineApi.CreatePipelineNotice)
		PipelinesRouter.POST("pipelines/closeCache/:id", PipelineApi.ClosePipelineCache)
		PipelinesRouter.POST("pipelines/cache", PipelineApi.CreatePipelineCache)
		PipelinesRouter.GET("pipelines/cache/:id", PipelineApi.GetPipelinesCache)
		PipelinesRouter.POST("pipelines/syncBranch/:id", PipelineApi.SyncBranches)
		PipelinesRouter.GET("pipelines/:id/branches", PipelineApi.GetBranchesList)

	}

}
