package organization

import (
	"DYCLOUD/global"
	"DYCLOUD/model/system"
	"DYCLOUD/plugin/organization/model"
	"DYCLOUD/plugin/organization/router"
	"DYCLOUD/plugin/plugin-tool/utils"
	"github.com/gin-gonic/gin"
)

type OrganizationPlugin struct {
}

func CreateOrganizationPlug() *OrganizationPlugin {
	global.DYCLOUD_DB.AutoMigrate(model.Organization{}, model.OrgUser{}, model.DataAuthority{})
	utils.RegisterMenus(
		system.SysBaseMenu{
			Path:      "organizationGroup",
			Name:      "organizationGroup",
			Hidden:    false,
			Component: "view/routerHolder.vue",
			Sort:      1000,
			Meta: system.Meta{
				Title: "组织管理",
				Icon:  "school",
			},
		},
		system.SysBaseMenu{
			Path:      "organization",
			Name:      "organization",
			Hidden:    false,
			Component: "plugin/organization/view/index.vue",
			Sort:      0,
			Meta: system.Meta{
				Title: "组织管理",
				Icon:  "school",
			},
		},
		system.SysBaseMenu{
			Path:      "dataAuthority",
			Name:      "dataAuthority",
			Hidden:    false,
			Component: "plugin/organization/view/dataAuthority.vue",
			Sort:      1,
			Meta: system.Meta{
				Title: "资源管理",
				Icon:  "money",
			},
		},
	)
	utils.RegisterApis(
		system.SysApi{
			Path:        "/org/createOrganization",
			Description: "创建组织",
			ApiGroup:    "组织管理",
			Method:      "POST",
		},
		system.SysApi{
			Path:        "/org/deleteOrganization",
			Description: "删除组织",
			ApiGroup:    "组织管理",
			Method:      "DELETE",
		},
		system.SysApi{
			Path:        "/org/deleteOrganizationByIds",
			Description: "批量删除组织",
			ApiGroup:    "组织管理",
			Method:      "DELETE",
		},
		system.SysApi{
			Path:        "/org/updateOrganization",
			Description: "更新组织",
			ApiGroup:    "组织管理",
			Method:      "PUT",
		},
		system.SysApi{
			Path:        "/org/createOrgUser",
			Description: "添加组织成员",
			ApiGroup:    "组织管理",
			Method:      "POST",
		},
		system.SysApi{
			Path:        "/org/setOrgUserAdmin",
			Description: "设置管理员",
			ApiGroup:    "组织管理",
			Method:      "PUT",
		},
		system.SysApi{
			Path:        "/org/findOrganization",
			Description: "获取所选组织",
			ApiGroup:    "组织管理",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/org/getOrganizationList",
			Description: "获取组织列表",
			ApiGroup:    "组织管理",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/org/getOrganizationTree",
			Description: "获取所有组织树",
			ApiGroup:    "组织管理",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/org/findOrgUserAll",
			Description: "获取组织下所有用户",
			ApiGroup:    "组织管理",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/org/findOrgUserList",
			Description: "获取组织下所有用户（分页）",
			ApiGroup:    "组织管理",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/org/deleteOrgUser",
			Description: "删除当前组织下选中用户",
			ApiGroup:    "组织管理",
			Method:      "DELETE",
		},
		system.SysApi{
			Path:        "/org/transferOrgUser",
			Description: "用户转移组织",
			ApiGroup:    "组织管理",
			Method:      "PUT",
		},
		system.SysApi{
			Path:        "/org/setDataAuthority",
			Description: "设置资源权限",
			ApiGroup:    "组织管理",
			Method:      "PUT",
		},
		system.SysApi{
			Path:        "/org/syncAuthority",
			Description: "同步角色数据",
			ApiGroup:    "组织管理",
			Method:      "POST",
		},
		system.SysApi{
			Path:        "/org/getAuthority",
			Description: "获取所有资源权限",
			ApiGroup:    "组织管理",
			Method:      "GET",
		},
	)
	return &OrganizationPlugin{}
}

func (*OrganizationPlugin) Register(group *gin.RouterGroup) {
	router.RouterGroupApp.InitOrganizationRouter(group)
}

func (*OrganizationPlugin) RouterPath() string {
	return "org"
}
