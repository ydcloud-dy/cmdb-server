package cloudtty

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cloudTTY"
	"DYCLOUD/model/common/response"
	"DYCLOUD/service"
	sutils "DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type K8sCloudTTYApi struct{}

var k8sCloudTTYService = service.ServiceGroupApp.CloudTTYServiceGroup.K8sCloudTTYService

func (t *K8sCloudTTYApi) CloudTTYGet(c *gin.Context) {
	var tty cloudTTY.CloudTTY
	_ = c.ShouldBindJSON(&tty)

	podMsg, err := k8sCloudTTYService.Get(tty, sutils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
	} else {
		response.OkWithDetailed(podMsg, "获取成功", c)
	}
}
