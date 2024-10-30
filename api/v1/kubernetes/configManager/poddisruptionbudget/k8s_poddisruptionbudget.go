package poddisruptionbudget

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/model/kubernetes/poddisruptionbudget"
	"DYCLOUD/service"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type K8sPoddisruptionbudgetApi struct{}

var k8sPoddisruptionbudgetService = service.ServiceGroupApp.PoddistruptionbudgetServiceGroup.K8sPoddisruptionbudgetService

func (k *K8sPoddisruptionbudgetApi) GetPoddisruptionbudgetList(c *gin.Context) {
	req := poddisruptionbudget.GetPoddisruptionbudgetListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := k8sPoddisruptionbudgetService.GetPoddisruptionbudgetList(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	} else {

		response.OkWithDetailed(poddisruptionbudget.PoddisruptionbudgetListResponse{
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
func (k *K8sPoddisruptionbudgetApi) DescribePoddisruptionbudgetInfo(c *gin.Context) {
	req := poddisruptionbudget.DescribePoddisruptionbudgetReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sPoddisruptionbudgetService.DescribePoddisruptionbudget(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(poddisruptionbudget.DescribePoddisruptionbudgetResponse{Items: list}, "获取成功", c)
	}
}
func (k *K8sPoddisruptionbudgetApi) UpdatePoddisruptionbudget(c *gin.Context) {
	req := poddisruptionbudget.UpdatePoddisruptionbudgetReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sPoddisruptionbudgetService.UpdatePoddisruptionbudget(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}
func (k *K8sPoddisruptionbudgetApi) DeletePoddisruptionbudget(c *gin.Context) {
	req := poddisruptionbudget.DeletePoddisruptionbudgetReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := k8sPoddisruptionbudgetService.DeletePoddisruptionbudget(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	} else {
		time.Sleep(1 * time.Second)
		response.OkWithMessage("删除成功", c)
	}
}
func (k *K8sPoddisruptionbudgetApi) CreatePoddisruptionbudget(c *gin.Context) {
	req := poddisruptionbudget.CreatePoddisruptionbudgetReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if CronJob, err := k8sPoddisruptionbudgetService.CreatePoddisruptionbudget(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(CronJob, c)
	}
}
