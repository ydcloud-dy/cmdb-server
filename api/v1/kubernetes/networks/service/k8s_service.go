package service

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	service2 "DYCLOUD/model/kubernetes/service"
	"DYCLOUD/service"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type K8sServiceApi struct{}

var k8sServiceService = service.ServiceGroupApp.SvcServiceGroup.K8sSvcService

func (k *K8sServiceApi) GetServiceList(c *gin.Context) {
	req := service2.GetServiceListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := k8sServiceService.GetServiceList(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(service2.ServiceListResponse{
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
func (k *K8sServiceApi) DescribeServiceInfo(c *gin.Context) {
	req := service2.DescribeServiceReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sServiceService.DescribeService(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(service2.DescribeServiceResponse{Items: list}, "获取成功", c)
	}
}
func (k *K8sServiceApi) UpdateService(c *gin.Context) {
	req := service2.UpdateServiceReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sServiceService.UpdateService(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}
func (k *K8sServiceApi) DeleteService(c *gin.Context) {
	req := service2.DeleteServiceReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := k8sServiceService.DeleteService(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	} else {
		time.Sleep(1 * time.Second)
		response.OkWithMessage("删除成功", c)
	}
}
func (k *K8sServiceApi) CreateService(c *gin.Context) {
	req := service2.CreateServiceReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if CronJob, err := k8sServiceService.CreateService(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(CronJob, c)
	}
}
