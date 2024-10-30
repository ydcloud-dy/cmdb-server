package docker

import (
	gvaGlobal "DYCLOUD/global"
	"DYCLOUD/model/docker"
	"DYCLOUD/model/system"
	"DYCLOUD/plugin/plugin-tool/utils"
	"DYCLOUD/router/docker"
	"github.com/gin-gonic/gin"
)

type Docker struct {
}

func CreateDockerPlug() *Docker {

	gvaGlobal.DYCLOUD_DB.AutoMigrate(model.Host{}) // 此处可以把插件依赖的数据库结构体自动创建表 需要填写对应的结构体

	// 具体值请根据实际情况修改
	utils.RegisterMenus(
		system.SysBaseMenu{
			Path:      "Docker",
			Name:      "Docker",
			Hidden:    false,
			Component: "view/routerHolder.vue",
			Sort:      0,
			Meta: system.Meta{
				Title: "Docker",
				Icon:  "monitor",
			},
		},
		system.SysBaseMenu{
			Path:      "hs",
			Name:      "hs",
			Hidden:    false,
			Component: "plugin/docker/view/host/host.vue",
			Sort:      0,
			Meta: system.Meta{
				Title: "主机列表",
				Icon:  "list",
			},
		},
		system.SysBaseMenu{
			Path:      "DockerDashboard",
			Name:      "DockerDashboard",
			Hidden:    false,
			Component: "plugin/docker/view/dashboard/index.vue",
			Sort:      1,
			Meta: system.Meta{
				Title: "仪表盘",
				Icon:  "aim",
			},
		},
		system.SysBaseMenu{
			Path:      "Terminal",
			Name:      "Terminal",
			Hidden:    true,
			Component: "plugin/docker/view/container/terminal.vue",
			Sort:      1,
			Meta: system.Meta{
				Title: "终端",
				Icon:  "operation",
			},
		},
		system.SysBaseMenu{
			Path:      "Container",
			Name:      "Container",
			Hidden:    false,
			Component: "plugin/docker/view/container/index.vue",
			Sort:      2,
			Meta: system.Meta{
				Title: "容器",
				Icon:  "grid",
			},
		},
		system.SysBaseMenu{
			Path:      "Image",
			Name:      "Image",
			Hidden:    false,
			Component: "plugin/docker/view/image/index.vue",
			Sort:      3,
			Meta: system.Meta{
				Title: "镜像",
				Icon:  "expand",
			},
		},
		system.SysBaseMenu{
			Path:      "Network",
			Name:      "Network",
			Hidden:    false,
			Component: "plugin/docker/view/network/index.vue",
			Sort:      4,
			Meta: system.Meta{
				Title: "网络",
				Icon:  "rank",
			},
		},
		system.SysBaseMenu{
			Path:      "Volume",
			Name:      "Volume",
			Hidden:    false,
			Component: "plugin/docker/view/volume/index.vue",
			Sort:      5,
			Meta: system.Meta{
				Title: "存储卷",
				Icon:  "coin",
			},
		},
	)

	// 下方会自动注册api 以下格式为示例格式，请按照实际情况修改
	utils.RegisterApis(
		system.SysApi{
			Path:        "/docker/host/genTlsScript",
			Description: "获取host生成tls脚本",
			ApiGroup:    "docker",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/docker/:host/container/inspect",
			Description: "获取容器详细信息",
			ApiGroup:    "docker",
			Method:      "POST",
		},
		system.SysApi{
			Path:        "/docker/host/getHostList",
			Description: "获取主机列表列表",
			ApiGroup:    "docker",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/docker/host/findHost",
			Description: "根据ID获取主机列表",
			ApiGroup:    "docker",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/docker/host/updateHost",
			Description: "更新主机列表",
			ApiGroup:    "docker",
			Method:      "PUT",
		}, system.SysApi{
			Path:        "/docker/host/deleteHostByIds",
			Description: "批量删除主机列表",
			ApiGroup:    "docker",
			Method:      "DELETE",
		},
		system.SysApi{
			Path:        "/docker/host/deleteHost",
			Description: "删除主机列表",
			ApiGroup:    "docker",
			Method:      "DELETE",
		},
		system.SysApi{
			Path:        "/docker/host/createHost",
			Description: "新增主机列表",
			ApiGroup:    "docker",
			Method:      "POST",
		},
		system.SysApi{
			Path:        "/docker/:host/volume",
			Description: "创建存储卷",
			ApiGroup:    "docker",
			Method:      "POST",
		},
		system.SysApi{
			Path:        "/docker/:host/system/info",
			Description: "系统信息",
			ApiGroup:    "docker",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/docker/:host/volume",
			Description: "删除存储卷",
			ApiGroup:    "docker",
			Method:      "DELETE",
		}, system.SysApi{
			Path:        "/docker/:host/volume/list",
			Description: "获取存储卷列表",
			ApiGroup:    "docker",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/docker/:host/network",
			Description: "创建网络",
			ApiGroup:    "docker",
			Method:      "POST",
		},
		system.SysApi{
			Path:        "/docker/:host/network",
			Description: "删除网络",
			ApiGroup:    "docker",
			Method:      "DELETE",
		},
		system.SysApi{
			Path:        "/docker/:host/network/list",
			Description: "获取网络列表",
			ApiGroup:    "docker",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/docker/:host/image",
			Description: "删除镜像",
			ApiGroup:    "docker",
			Method:      "DELETE",
		},
		system.SysApi{
			Path:        "/docker/:host/image/pull",
			Description: "下载镜像",
			ApiGroup:    "docker",
			Method:      "POST",
		},
		system.SysApi{
			Path:        "/docker/:host/image/list",
			Description: "获取镜像列表",
			ApiGroup:    "docker",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/docker/:host/container/exec/term/:id",
			Description: "连接容器终端",
			ApiGroup:    "docker",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/docker/:host/container/exec/resize",
			Description: "修改容器终端大小",
			ApiGroup:    "docker",
			Method:      "POST",
		},
		system.SysApi{
			Path:        "/docker/:host/container/exec",
			Description: "创建容器终端",
			ApiGroup:    "docker",
			Method:      "POST",
		},
		system.SysApi{
			Path:        "/docker/:host/container/stats",
			Description: "获取容器实时监控数据",
			ApiGroup:    "docker",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/docker/:host/container/start",
			Description: "启动容器",
			ApiGroup:    "docker",
			Method:      "POST",
		},
		system.SysApi{
			Path:        "/docker/:host/container/stop",
			Description: "停止容器",
			ApiGroup:    "docker",
			Method:      "POST",
		},
		system.SysApi{
			Path:        "/docker/:host/container/restart",
			Description: "重启容器",
			ApiGroup:    "docker",
			Method:      "POST",
		},
		system.SysApi{
			Path:        "/docker/:host/container",
			Description: "创建容器",
			ApiGroup:    "docker",
			Method:      "POST",
		},
		system.SysApi{
			Path:        "/docker/:host/container",
			Description: "删除容器",
			ApiGroup:    "docker",
			Method:      "DELETE",
		},
		system.SysApi{
			Path:        "/docker/:host/container/logs",
			Description: "获取容器日志",
			ApiGroup:    "docker",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/docker/:host/container/list",
			Description: "获取容器列表",
			ApiGroup:    "docker",
			Method:      "GET",
		},
		system.SysApi{
			Path:        "/docker/host/checkHost",
			Description: "检测主机连接可用性",
			ApiGroup:    "docker",
			Method:      "POST",
		},
	)

	return &Docker{}
}

func (*Docker) Register(group *gin.RouterGroup) {
	router.RouterGroupApp.InitDockerRouter(group)
}

func (*Docker) RouterPath() string {
	return "docker"
}
