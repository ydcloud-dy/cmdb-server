package cronjob

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	cronjob2 "DYCLOUD/model/kubernetes/cronjob"
	"DYCLOUD/service"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type K8sCronJobApi struct{}

var k8sCronCronJobService = service.ServiceGroupApp.CronJobServiceGroup

// GetCronJobList 获取cronjob列表
// @Tags kubernetes
// @Summary 获取cronjob列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "获取cronjob列表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /kubernetes/cronJob [get]
func (k *K8sCronJobApi) GetCronJobList(c *gin.Context) {
	req := cronjob2.GetCronJobListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := k8sCronCronJobService.GetCronJobList(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	} else {

		response.OkWithDetailed(cronjob2.CronJobListResponse{
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

// DescribeCronJobInfo 获取 CronJob 详情
// @Tags kubernetes
// @Summary 获取 CronJob 详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cronjob2.DescribeCronJobReq true "获取 CronJob 详情"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /kubernetes/cronJobDetails [get]
func (k *K8sCronJobApi) DescribeCronJobInfo(c *gin.Context) {
	req := cronjob2.DescribeCronJobReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sCronCronJobService.DescribeCronJob(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(cronjob2.DescribeCronJobResponse{Items: list}, "获取成功", c)
	}
}

// UpdateCronJob 更新 CronJob
// @Tags kubernetes
// @Summary 更新 CronJob
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cronjob2.UpdateCronJobReq true "更新 CronJob"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /kubernetes/cronJob [put]
func (k *K8sCronJobApi) UpdateCronJob(c *gin.Context) {
	req := cronjob2.UpdateCronJobReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sCronCronJobService.UpdateCronJob(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}

// DeleteCronJob 删除 CronJob
// @Tags kubernetes
// @Summary 删除 CronJob
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cronjob2.DeleteCronJobReq true "删除 CronJob"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /kubernetes/cronJob [delete]
func (k *K8sCronJobApi) DeleteCronJob(c *gin.Context) {
	req := cronjob2.DeleteCronJobReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := k8sCronCronJobService.DeleteCronJob(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	} else {
		time.Sleep(1 * time.Second)
		response.OkWithMessage("删除成功", c)
	}
}

// CreateCronJob 创建 CronJob
// @Tags kubernetes
// @Summary 创建 CronJob
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cronjob2.CreateCronJobReq true "创建 CronJob"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /kubernetes/cronJob [post]
func (k *K8sCronJobApi) CreateCronJob(c *gin.Context) {
	req := cronjob2.CreateCronJobReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if CronJob, err := k8sCronCronJobService.CreateCronJob(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(CronJob, c)
	}
}
