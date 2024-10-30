package initialize

import (
	"DYCLOUD/global"
	"DYCLOUD/middleware/cloudCmdb"
	"DYCLOUD/middleware/docker"
	"DYCLOUD/plugin/email"
	"DYCLOUD/plugin/organization"
	"DYCLOUD/utils/plugin"
	"fmt"
	"github.com/gin-gonic/gin"
)

func PluginInit(group *gin.RouterGroup, Plugin ...plugin.Plugin) {
	for i := range Plugin {
		fmt.Println(Plugin[i].RouterPath(), "注册开始!")
		PluginGroup := group.Group(Plugin[i].RouterPath())
		Plugin[i].Register(PluginGroup)
		fmt.Println(Plugin[i].RouterPath(), "注册成功!")
	}
}

func bizPluginV1(group ...*gin.RouterGroup) {
	private := group[0]
	public := group[1]
	//  添加跟角色挂钩权限的插件 示例 本地示例模式于在线仓库模式注意上方的import 可以自行切换 效果相同
	PluginInit(private, email.CreateEmailPlug(
		global.DYCLOUD_CONFIG.Email.To,
		global.DYCLOUD_CONFIG.Email.From,
		global.DYCLOUD_CONFIG.Email.Host,
		global.DYCLOUD_CONFIG.Email.Secret,
		global.DYCLOUD_CONFIG.Email.Nickname,
		global.DYCLOUD_CONFIG.Email.Port,
		global.DYCLOUD_CONFIG.Email.IsSSL,
	))
	PluginInit(private, docker.CreateDockerPlug())
	PluginInit(public, organization.CreateOrganizationPlug())
	PluginInit(private, cloudCmdb.CreateCloudCmdbPlug())
	holder(public, private)
}
