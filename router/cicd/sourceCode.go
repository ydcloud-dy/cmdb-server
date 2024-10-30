package cicd

import (
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type SourceCodeRouter struct{}

// InitServiceIntrgrationRouter 初始化 ServiceIntrgration表 路由信息
func (s *SourceCodeRouter) InitSourceCodeRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	SourceCodeRouter := Router.Group("cicd").Use(middleware.OperationRecord())
	//ServiceIntrgrationRouterWithoutRecord := Router.Group("cmdb")
	//ServiceIntrgrationRouterWithoutAuth := PublicRouter.Group("cmdb")
	{
		SourceCodeRouter.GET("sourceCode", SourceCodeApi.GetSourceCodeList)
		SourceCodeRouter.POST("sourceCode", SourceCodeApi.CreateSourceCode)
		SourceCodeRouter.PUT("sourceCode", SourceCodeApi.UpdateSourceCode)
		SourceCodeRouter.DELETE("sourceCode/:id", SourceCodeApi.DeleteSourceCode)
		SourceCodeRouter.GET("sourceCode/:id", SourceCodeApi.DescribeSourceCode)
		SourceCodeRouter.POST("sourceCode/verify", SourceCodeApi.VerifySourceCode)

	}

}
