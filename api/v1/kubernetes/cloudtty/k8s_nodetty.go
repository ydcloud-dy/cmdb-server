package cloudtty

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/response"
	kubernetesReq "DYCLOUD/model/kubernetes/cluster/request"
	"DYCLOUD/service"
	sutils "DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type K8SNodeTTYApi struct{}

var k8sNodeTTYService = service.ServiceGroupApp.CloudTTYServiceGroup.NodeTTYService

// @Tags NodeTTyApi
// @Summary  NodeTTY 终端资源创建
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Router  /kubernetes/nodetty/get [POST]
func (n *K8SNodeTTYApi) NodeTTYGet(c *gin.Context) {
	var nodetty kubernetesReq.NodeTTY
	_ = c.ShouldBindJSON(&nodetty)

	podMsg, err := k8sNodeTTYService.Get(nodetty, sutils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
	} else {
		response.OkWithDetailed(podMsg, "获取成功", c)
	}
}
