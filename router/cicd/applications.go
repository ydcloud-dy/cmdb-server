package cicd

import (
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type ApplicationsRouter struct{}

// InitApplicationsRouter 初始化 Applications表 路由信息
func (s *ApplicationsRouter) InitApplicationsRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	ApplicationsRouter := Router.Group("cicd").Use(middleware.OperationRecord())
	{
		ApplicationsRouter.GET("applications", ApplicationApi.GetApplicationsList)
		ApplicationsRouter.POST("applications", ApplicationApi.CreateApplications)
		ApplicationsRouter.PUT("applications", ApplicationApi.UpdateApplications)
		ApplicationsRouter.DELETE("applications/:id", ApplicationApi.DeleteApplications)
		ApplicationsRouter.GET("applications/:id", ApplicationApi.DescribeApplications)
		ApplicationsRouter.GET("applicationsByName", ApplicationApi.DescribeApplicationsByName)

		ApplicationsRouter.DELETE("applications", ApplicationApi.DeleteApplicationsByIds)
		ApplicationsRouter.POST("applications/:id/syncBranches", ApplicationApi.SyncBranches)
		ApplicationsRouter.GET("applications/:id/branches", ApplicationApi.GetBranchesList)
		ApplicationsRouter.POST("applications/deploymentInfo", ApplicationApi.GetApplicationDeploymentInfo)
	}

}
