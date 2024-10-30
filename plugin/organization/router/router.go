package router

import (
	"DYCLOUD/middleware"
	"DYCLOUD/plugin/organization/api"
	"github.com/gin-gonic/gin"
)

type OrganizationRouter struct{}

func (s *OrganizationRouter) InitOrganizationRouter(Router *gin.RouterGroup) {
	orgRouter := Router.Use(middleware.OperationRecord())
	orgRouterWithoutRecord := Router
	var orgApi = api.ApiGroupApp.OrganizationApi
	{
		orgRouter.POST("createOrganization", orgApi.CreateOrganization)             // 新建Organization
		orgRouter.DELETE("deleteOrganization", orgApi.DeleteOrganization)           // 删除Organization
		orgRouter.DELETE("deleteOrganizationByIds", orgApi.DeleteOrganizationByIds) // 批量删除Organization
		orgRouter.PUT("updateOrganization", orgApi.UpdateOrganization)              // 更新Organization
		orgRouter.POST("createOrgUser", orgApi.CreateOrgUser)                       // 人员入职
		orgRouter.PUT("setOrgUserAdmin", orgApi.SetOrgUserAdmin)                    // 管理员设置
		orgRouter.PUT("setDataAuthority", orgApi.SetOrgAuthority)                   // 设置资源权限
		orgRouter.POST("syncAuthority", orgApi.SyncAuthority)                       // 同步角色
		orgRouter.GET("getAuthority", orgApi.GetAuthority)                          // 获取资源权限
	}
	{
		orgRouterWithoutRecord.GET("findOrganization", orgApi.FindOrganization)       // 根据ID获取Organization
		orgRouterWithoutRecord.GET("getOrganizationList", orgApi.GetOrganizationList) // 获取Organization列表
		orgRouterWithoutRecord.GET("getOrganizationTree", orgApi.GetOrganizationTree) // 获取Organization树
		orgRouterWithoutRecord.GET("findOrgUserAll", orgApi.FindOrgUserAll)           // 获取当前组织下所有用户ID
		orgRouterWithoutRecord.GET("findOrgUserList", orgApi.FindOrgUserList)         // 获取当前组织下所有用户(分页)
		orgRouterWithoutRecord.DELETE("deleteOrgUser", orgApi.DeleteOrgUser)          // 删除当前组织下选中用户
		orgRouterWithoutRecord.PUT("transferOrgUser", orgApi.TransferOrgUser)         // 用户转移组织
	}
}
