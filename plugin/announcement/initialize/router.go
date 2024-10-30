package initialize

import (
	"DYCLOUD/global"
	"DYCLOUD/middleware"
	"DYCLOUD/plugin/announcement/router"
	"github.com/gin-gonic/gin"
)

func Router(engine *gin.Engine) {
	public := engine.Group(global.DYCLOUD_CONFIG.System.RouterPrefix).Group("")
	private := engine.Group(global.DYCLOUD_CONFIG.System.RouterPrefix).Group("")
	private.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	router.Router.Info.Init(public, private)
}
