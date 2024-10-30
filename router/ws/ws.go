package ws

import (
	api "DYCLOUD/api/v1"
	"github.com/gin-gonic/gin"
)

type WsApiRouter struct{}

func (t *WsApiRouter) InitWsRouter(Router *gin.RouterGroup) {
	WsRouter := Router.Group("kubernetes/pods")
	wsApi := api.ApiGroupApp.WsApi
	podsApi := api.ApiGroupApp.PodApiGroup
	WsRouter.GET("/terminal", wsApi.Terminal)           // 终端
	WsRouter.GET("/logs", wsApi.ContainerLog)           // 终端日志
	WsRouter.GET("/downloadFile", podsApi.DownloadFile) // 文件下载
}
