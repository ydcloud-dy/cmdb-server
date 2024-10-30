package cicd

import (
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type EnvironmentRouter struct{}

// InitEnvironmentRouter 初始化 Environment表 路由信息
func (s *EnvironmentRouter) InitEnvironmentRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	EnvironmentRouter := Router.Group("cicd").Use(middleware.OperationRecord())
	//EnvironmentRouterWithoutRecord := Router.Group("cmdb")
	//EnvironmentRouterWithoutAuth := PublicRouter.Group("cmdb")
	{
		EnvironmentRouter.GET("environment", EnvironmentApi.GetEnvironmentList)
		EnvironmentRouter.POST("environment", EnvironmentApi.CreateEnvironment)
		EnvironmentRouter.PUT("environment", EnvironmentApi.UpdateEnvironment)
		EnvironmentRouter.DELETE("environment/:id", EnvironmentApi.DeleteEnvironment)
		EnvironmentRouter.GET("environment/:id", EnvironmentApi.DescribeEnvironment)
		EnvironmentRouter.DELETE("environment", EnvironmentApi.DeleteEnvironmentByIds)
	}
	{
		//EnvironmentRouterWithoutRecord.GET("hostsById", EnvironmentApi.FindEnvironment) // 根据ID获取Environment表
		//EnvironmentRouterWithoutRecord.GET("hosts", EnvironmentApi.GetEnvironmentList)  // 获取Environment表列表
	}
	{
		//EnvironmentRouterWithoutAuth.GET("hostsPublic", EnvironmentApi.GetEnvironmentPublic) // Environment表开放接口
	}
}
