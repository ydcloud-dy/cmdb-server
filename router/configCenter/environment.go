package configCenter

import (
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type EnvironmentRouter struct{}

// InitEnvironmentRouter 初始化 Environment表 路由信息
func (s *EnvironmentRouter) InitEnvironmentRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	EnvironmentRouter := Router.Group("configCenter").Use(middleware.OperationRecord())
	{
		EnvironmentRouter.GET("environment", EnvironmentApi.GetEnvironmentList)
		EnvironmentRouter.POST("environment", EnvironmentApi.CreateEnvironment)
		EnvironmentRouter.PUT("environment", EnvironmentApi.UpdateEnvironment)
		EnvironmentRouter.DELETE("environment/:id", EnvironmentApi.DeleteEnvironment)
		EnvironmentRouter.GET("environment/:id", EnvironmentApi.DescribeEnvironment)
		EnvironmentRouter.DELETE("environment", EnvironmentApi.DeleteEnvironmentByIds)
	}

}
