package configCenter

import (
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type ApplicationsRouter struct{}

// InitApplicationsRouter 初始化 Applications表 路由信息
func (s *ApplicationsRouter) InitApplicationsRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	ApplicationsRouter := Router.Group("configCenter").Use(middleware.OperationRecord())
	{
		ApplicationsRouter.GET("Applications", ApplicationApi.GetApplicationsList)
		ApplicationsRouter.POST("Applications", ApplicationApi.CreateApplications)
		ApplicationsRouter.PUT("Applications", ApplicationApi.UpdateApplications)
		ApplicationsRouter.DELETE("Applications/:id", ApplicationApi.DeleteApplications)
		ApplicationsRouter.GET("Applications/:id", ApplicationApi.DescribeApplications)
		ApplicationsRouter.DELETE("Applications", ApplicationApi.DeleteApplicationsByIds)
	}

}
