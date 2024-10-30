package cloudCmdb

import (
	"DYCLOUD/global"
	model "DYCLOUD/model/cloudCmdb"
	sysModel "DYCLOUD/model/system"
	"DYCLOUD/plugin/plugin-tool/utils"
	router "DYCLOUD/router/cloudCmdb"
	"github.com/gin-gonic/gin"
)

type CloudCmdbPlugin struct {
}

func CreateCloudCmdbPlug() *CloudCmdbPlugin {
	global.DYCLOUD_DB.AutoMigrate(
		model.CloudPlatform{},
		model.CloudRegions{},
		model.VirtualMachine{},
		model.LoadBalancer{},
		model.RDS{})

	utils.RegisterApis(
		sysModel.SysApi{ApiGroup: "云厂商管理", Method: "POST", Path: "/cloudcmdb/cloud_platform/list", Description: "列表"},
		sysModel.SysApi{ApiGroup: "云厂商管理", Method: "POST", Path: "/cloudcmdb/cloud_platform/getById", Description: "单个获取"},
		sysModel.SysApi{ApiGroup: "云厂商管理", Method: "POST", Path: "/cloudcmdb/cloud_platform/create", Description: "创建"},
		sysModel.SysApi{ApiGroup: "云厂商管理", Method: "PUT", Path: "/cloudcmdb/cloud_platform/update", Description: "更新"},
		sysModel.SysApi{ApiGroup: "云厂商管理", Method: "DELETE", Path: "/cloudcmdb/cloud_platform/delete", Description: "删除"},
		sysModel.SysApi{ApiGroup: "云厂商管理", Method: "DELETE", Path: "/cloudcmdb/cloud_platform/deleteByIds", Description: "批量删除"},
		sysModel.SysApi{ApiGroup: "云厂商区域管理", Method: "POST", Path: "/cloudcmdb/cloud_region/syncRegion", Description: "区域同步"},
		sysModel.SysApi{ApiGroup: "云主机管理", Method: "POST", Path: "/cloudcmdb/virtualMachine/tree", Description: "目录树"},
		sysModel.SysApi{ApiGroup: "云主机管理", Method: "POST", Path: "/cloudcmdb/virtualMachine/list", Description: "列表"},
		sysModel.SysApi{ApiGroup: "云主机管理", Method: "POST", Path: "/cloudcmdb/virtualMachine/sync", Description: "同步"},
		sysModel.SysApi{ApiGroup: "负载均衡管理", Method: "POST", Path: "/cloudcmdb/loadBalancer/tree", Description: "目录树"},
		sysModel.SysApi{ApiGroup: "负载均衡管理", Method: "POST", Path: "/cloudcmdb/loadBalancer/list", Description: "列表"},
		sysModel.SysApi{ApiGroup: "负载均衡管理", Method: "POST", Path: "/cloudcmdb/loadBalancer/sync", Description: "同步"},
		sysModel.SysApi{ApiGroup: "云数据库管理", Method: "POST", Path: "/cloudcmdb/rds/tree", Description: "目录树"},
		sysModel.SysApi{ApiGroup: "云数据库管理", Method: "POST", Path: "/cloudcmdb/rds/list", Description: "列表"},
		sysModel.SysApi{ApiGroup: "云数据库管理", Method: "POST", Path: "/cloudcmdb/rds/sync", Description: "同步"},
	)

	utils.RegisterMenus(
		sysModel.SysBaseMenu{
			Name:      "资产管理",
			Path:      "cloudcmdb",
			Hidden:    false,
			Component: "plugin/cloudcmdb/view/index.vue",
			Sort:      1000,
			Meta:      sysModel.Meta{Title: "云资产管理", Icon: "cloudy"},
		},
		sysModel.SysBaseMenu{
			Name:      "cloud_platform",
			Path:      "cloud_platform",
			Hidden:    false,
			Component: "plugin/cloudcmdb/view/cloud_platform/index.vue",
			Sort:      0,
			Meta:      sysModel.Meta{Title: "云厂商", Icon: "cloudy"},
		},
		sysModel.SysBaseMenu{
			Name:      "cloud_load_balancer",
			Path:      "cloud_load_balancer",
			Hidden:    false,
			Component: "plugin/cloudcmdb/view/cloud_load_balancer/index.vue",
			Sort:      1,
			Meta:      sysModel.Meta{Title: "负载均衡", Icon: "Grid"},
		},
		sysModel.SysBaseMenu{
			Name:      "cloud_virtual_machine",
			Path:      "cloud_virtual_machine",
			Hidden:    false,
			Component: "plugin/cloudcmdb/view/cloud_virtual_machine/index.vue",
			Sort:      2,
			Meta:      sysModel.Meta{Title: "云服务器", Icon: "film"},
		},
		sysModel.SysBaseMenu{
			Name:      "cloud_rds",
			Path:      "cloud_rds",
			Hidden:    false,
			Component: "plugin/cloudcmdb/view/cloud_rds/index.vue",
			Sort:      3,
			Meta:      sysModel.Meta{Title: "云数据库", Icon: "coin"},
		},
	)

	return &CloudCmdbPlugin{}
}

func (*CloudCmdbPlugin) Register(group *gin.RouterGroup) {
	router.RouterGroupApp.InitCloudPlatformRouter(group)
	router.RouterGroupApp.InitCloudRegionRouter(group)
	router.RouterGroupApp.InitVirtualMachineRouter(group)
	router.RouterGroupApp.InitLoadBalancerRouter(group)
	router.RouterGroupApp.InitRDSRouter(group)
}

func (*CloudCmdbPlugin) RouterPath() string {
	return "cloudcmdb"
}
