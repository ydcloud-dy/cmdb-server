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
	}

}