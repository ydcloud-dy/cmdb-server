package configCenter

import (
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type ServiceIntegrationRouter struct{}

// InitServiceIntrgrationRouter 初始化 ServiceIntrgration表 路由信息
func (s *ServiceIntegrationRouter) InitServiceIntrgrationRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	ServiceIntegrationRouter := Router.Group("configCenter").Use(middleware.OperationRecord())
	//ServiceIntrgrationRouterWithoutRecord := Router.Group("cmdb")
	//ServiceIntrgrationRouterWithoutAuth := PublicRouter.Group("cmdb")
	{
		ServiceIntegrationRouter.GET("service", ServiceIntegrationApi.GetServiceIntegrationList)
		ServiceIntegrationRouter.GET("service/registry", ServiceIntegrationApi.GetRegistryList)

		ServiceIntegrationRouter.POST("service", ServiceIntegrationApi.CreateServiceIntegration)
		ServiceIntegrationRouter.PUT("service", ServiceIntegrationApi.UpdateServiceIntegration)
		ServiceIntegrationRouter.DELETE("service/:id", ServiceIntegrationApi.DeleteServiceIntegration)
		ServiceIntegrationRouter.GET("service/:id", ServiceIntegrationApi.DescribeServiceIntegration)
		ServiceIntegrationRouter.POST("service/verify", ServiceIntegrationApi.VerifyServiceIntegration)

	}

}
