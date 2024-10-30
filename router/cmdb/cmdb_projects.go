package cmdb

import (
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type CmdbProjectsRouter struct{}

// InitCmdbProjectsRouter 初始化 cmdbProjects表 路由信息
func (s *CmdbProjectsRouter) InitCmdbProjectsRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	cmdbProjectsRouter := Router.Group("cmdb").Use(middleware.OperationRecord())
	cmdbProjectsRouterWithoutRecord := Router.Group("cmdb")
	cmdbProjectsRouterWithoutAuth := PublicRouter.Group("cmdb")
	{
		cmdbProjectsRouter.POST("projects", cmdbProjectsApi.CreateCmdbProjects)             // 新建cmdbProjects表
		cmdbProjectsRouter.DELETE("projects", cmdbProjectsApi.DeleteCmdbProjects)           // 删除cmdbProjects表
		cmdbProjectsRouter.DELETE("projectsByIds", cmdbProjectsApi.DeleteCmdbProjectsByIds) // 批量删除cmdbProjects表
		cmdbProjectsRouter.PUT("projects", cmdbProjectsApi.UpdateCmdbProjects)              // 更新cmdbProjects表
	}
	{
		cmdbProjectsRouterWithoutRecord.GET("projectsById", cmdbProjectsApi.FindCmdbProjects) // 根据ID获取cmdbProjects表
		cmdbProjectsRouterWithoutRecord.GET("projects", cmdbProjectsApi.GetCmdbProjectsList)  // 获取cmdbProjects表列表
	}
	{
		cmdbProjectsRouterWithoutAuth.GET("projectsPublic", cmdbProjectsApi.GetCmdbProjectsPublic) // cmdbProjects表开放接口
	}
}
