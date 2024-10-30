package horizontalpodautoscalers

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/model/kubernetes/horizontalPod"
	"DYCLOUD/service"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type K8sHorizontalApi struct{}

var k8sHorizontalService = service.ServiceGroupApp.HorizontalPodServiceGroup.K8sHorizontalPodService

func (k *K8sHorizontalApi) GetHorizontalList(c *gin.Context) {
	req := horizontalPod.GetHorizontalPodListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := k8sHorizontalService.GetHorizontalPodList(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	} else {

		response.OkWithDetailed(horizontalPod.HorizontalPodResponse{
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
func (k *K8sHorizontalApi) DescribeHorizontalInfo(c *gin.Context) {
	req := horizontalPod.DescribeHorizontalPodReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sHorizontalService.DescribeHorizontalPod(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(horizontalPod.DescribeHorizontalPodResponse{Items: list}, "获取成功", c)
	}
}
func (k *K8sHorizontalApi) UpdateHorizontal(c *gin.Context) {
	req := horizontalPod.UpdateHorizontalPodReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sHorizontalService.UpdateHorizontalPod(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}
func (k *K8sHorizontalApi) DeleteHorizontal(c *gin.Context) {
	req := horizontalPod.DeleteHorizontalPodReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := k8sHorizontalService.DeleteHorizontalPod(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	} else {
		time.Sleep(1 * time.Second)
		response.OkWithMessage("删除成功", c)
	}
}
func (k *K8sHorizontalApi) CreateHorizontal(c *gin.Context) {
	req := horizontalPod.CreateHorizontalPodReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if CronJob, err := k8sHorizontalService.CreateHorizontalPod(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(CronJob, c)
	}
}
