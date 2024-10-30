package api

import (
	"DYCLOUD/global"
	global2 "DYCLOUD/global/docker"
	"DYCLOUD/model/common/response"
	service "DYCLOUD/service/docker"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemApi struct{}

// SystemInfo @Tags Docker
// @Summary  docker 系统信息
// @Produce  application/json
// @Success 200
// @Router /docker/system/info [get]
func (p *SystemApi) SystemInfo(c *gin.Context) {

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}

	if res, err := service.ServiceGroupApp.Info(host); err != nil {
		global.DYCLOUD_LOG.Error("获取系统信息失败!", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取系统信息失败:%v", err), c)
	} else {
		response.OkWithDetailed(res, "成功", c)
	}
}
