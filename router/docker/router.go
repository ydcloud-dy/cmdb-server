package router

import (
	"DYCLOUD/api/v1/docker"
	"DYCLOUD/middleware"
	"github.com/gin-gonic/gin"
)

type DockerRouter struct {
}

func (s *DockerRouter) InitDockerRouter(Router *gin.RouterGroup) {
	plugRouter := Router.Group(":host").Use(middleware.OperationRecord())
	plugRouterWithOut := Router.Use(middleware.OperationRecord())

	plugApi := api.ApiGroupApp.ContainerApi
	{
		// 容器管理
		plugRouter.GET("/container/list", plugApi.ListContainer)               // 获取容器列表
		plugRouter.GET("/container/logs", plugApi.LogsContainer)               // 获取容器日志
		plugRouter.DELETE("/container", plugApi.RemoveContainer)               // 删除容器
		plugRouter.POST("/container", plugApi.AddContainer)                    // 创建容器
		plugRouter.POST("/container/restart", plugApi.RestartContainer)        // 重启容器
		plugRouter.POST("/container/stop", plugApi.StopContainer)              // 停止容器
		plugRouter.POST("/container/start", plugApi.StartContainer)            // 启动容器
		plugRouter.GET("/container/stats", plugApi.StatsContainer)             // 获取容器实时监控数据
		plugRouter.POST("/container/exec", plugApi.ExecContainer)              // 执行容器
		plugRouter.POST("/container/exec/resize", plugApi.ExecContainerResize) // 容器执行终端resize
		plugRouter.PUT("/container", plugApi.EditContainer)                    // 编辑容器
		plugRouter.GET("/container/exec/term/:id", plugApi.ExecTermContainer)  // 容器执行终端
		plugRouter.POST("/container/inspect", plugApi.InspectContainer)        // 获取容器详细信息

	}

	plugImageApi := api.ApiGroupApp.ImageApi
	{
		// 镜像管理
		plugRouter.GET("/image/list", plugImageApi.ListImage)  // 获取镜像列表
		plugRouter.POST("/image/pull", plugImageApi.PullImage) // 下载镜像
		plugRouter.DELETE("/image", plugImageApi.RemoveImage)  // 删除镜像
	}

	plugNetworkApi := api.ApiGroupApp.NetworkApi
	{
		// 网络管理
		plugRouter.GET("/network/list", plugNetworkApi.ListNetwork) // 获取网络列表
		plugRouter.DELETE("/network", plugNetworkApi.RemoveNetwork) // 删除网络
		plugRouter.POST("/network", plugNetworkApi.CreateNetwork)   // 创建网络

	}

	plugVolumeApi := api.ApiGroupApp.VolumeApi
	{
		// 存储卷管理
		plugRouter.GET("/volume/list", plugVolumeApi.ListVolume) // 获取存储卷列表
		plugRouter.DELETE("/volume", plugVolumeApi.RemoveVolume) // 删除存储卷
		plugRouter.POST("/volume", plugVolumeApi.CreateVolume)   // 创建存储卷

	}

	plugSystemApi := api.ApiGroupApp.SystemApi
	{
		// 系统信息
		plugRouter.GET("/system/info", plugSystemApi.SystemInfo) // 获取系统信息

	}

	var hsApi = api.ApiGroupApp.HostApi
	{
		plugRouterWithOut.POST("/host/createHost", hsApi.CreateHost)             // 新建主机列表
		plugRouterWithOut.DELETE("/host/deleteHost", hsApi.DeleteHost)           // 删除主机列表
		plugRouterWithOut.DELETE("/host/deleteHostByIds", hsApi.DeleteHostByIds) // 批量删除主机列表
		plugRouterWithOut.PUT("/host/updateHost", hsApi.UpdateHost)              // 更新主机列表
		plugRouterWithOut.GET("/host/findHost", hsApi.FindHost)                  // 根据ID获取主机列表
		plugRouterWithOut.GET("/host/getHostList", hsApi.GetHostList)            // 获取主机列表列表
		plugRouterWithOut.GET("/host/genTlsScript", hsApi.GetHostGenTlsScript)   // 获取host生成tls脚本
		plugRouterWithOut.POST("/host/checkHost", hsApi.CheckHost)               // 检测主机连接可用性

	}

}
