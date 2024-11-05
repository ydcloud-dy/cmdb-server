package configCenter

import (
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type BuildEnvRouter struct{}

// InitBuildEnvRouter 初始化 BuildEnv表 路由信息
func (s *BuildEnvRouter) InitBuildEnvRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	BuildEnvRouter := Router.Group("configCenter").Use(middleware.OperationRecord())
	{
		BuildEnvRouter.GET("buildEnv", BuildEnvApi.GetBuildEnvList)
		BuildEnvRouter.POST("buildEnv", BuildEnvApi.CreateBuildEnv)
		BuildEnvRouter.PUT("buildEnv", BuildEnvApi.UpdateBuildEnv)
		BuildEnvRouter.DELETE("buildEnv/:id", BuildEnvApi.DeleteBuildEnv)
		BuildEnvRouter.GET("buildEnv/:id", BuildEnvApi.DescribeBuildEnv)
		BuildEnvRouter.DELETE("buildEnv", BuildEnvApi.DeleteBuildEnvByIds)
	}

}
