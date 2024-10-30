package daemonSet

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/model/kubernetes/daemonset"
	"DYCLOUD/service"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type K8sDaemonSetApi struct{}

var k8sDaemonSetService = service.ServiceGroupApp.DaemonSetServiceGroup.K8sDaemonSetService

func (k *K8sDaemonSetApi) GetDaemonSetList(c *gin.Context) {
	req := daemonset.GetDaemonSetListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := k8sDaemonSetService.GetDaemonSetList(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(daemonset.DaemonSetListResponse{
			Items: list,
			Total: total,
			PageInfo: request.PageInfo{
				Page:     req.Page,
				PageSize: req.PageSize,
				Keyword:  req.Keyword,
			},
		}, "获取成功", c)
	}
}
func (k *K8sDaemonSetApi) DescribeDaemonSetInfo(c *gin.Context) {
	req := daemonset.DescribeDaemonSetReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sDaemonSetService.DescribeDaemonSet(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(daemonset.DescribeDaemonSetResponse{Items: list}, "获取成功", c)
	}
}
func (k *K8sDaemonSetApi) UpdateDaemonsSet(c *gin.Context) {
	req := daemonset.UpdateDaemonSetReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sDaemonSetService.UpdateDaemonSet(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}
func (k *K8sDaemonSetApi) DeleteDaemonSet(c *gin.Context) {
	req := daemonset.DeleteDaemonSetReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := k8sDaemonSetService.DeleteDaemonSet(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	} else {
		response.OkWithMessage("删除成功", c)
	}
}
func (k *K8sDaemonSetApi) CreateDaemonSet(c *gin.Context) {
	req := daemonset.CreateDaemonSetReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if daemonset, err := k8sDaemonSetService.CreateDaemonSet(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(daemonset, c)
	}
}
