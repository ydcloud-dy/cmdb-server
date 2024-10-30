package configmap

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/model/kubernetes/configmap"
	"DYCLOUD/service"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type K8sConfigMapApi struct{}

var k8sConfigMapService = service.ServiceGroupApp.ConfigMapServiceGroup.K8sConfigMapService

// GetConfigMapList 获取configmap列表
// @Tags kubernetes
// @Summary 获取configmap列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "获取configmap列表"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /kubernetes/configMap [get]
func (k K8sConfigMapApi) GetConfigMapList(c *gin.Context) {
	req := configmap.GetConfigMapListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := k8sConfigMapService.GetConfigMapList(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(configmap.ConfigMapListResponse{
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

// DescribeConfigMapInfo 获取configmap详情
// @Tags kubernetes
// @Summary 获取configmap详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "获取configmap详情"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /kubernetes/configMapDetails [get]
func (k *K8sConfigMapApi) DescribeConfigMapInfo(c *gin.Context) {
	req := configmap.DescribeConfigMapReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, err := k8sConfigMapService.DescribeConfigMap(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(configmap.DescribeConfigMapResponse{Items: list}, "获取成功", c)
	}
}

// UpdateConfigMap 更新configmap
// @Tags kubernetes
// @Summary 更新configmap
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "更新configmap"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /kubernetes/configMap [get]
func (k *K8sConfigMapApi) UpdateConfigMap(c *gin.Context) {
	req := configmap.UpdateConfigMapReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, err := k8sConfigMapService.UpdateConfigMap(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}

// DeleteConfigMap 删除configmap
// @Tags kubernetes
// @Summary 删除configmap
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "删除configmap"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /kubernetes/configMap [delete]
func (k *K8sConfigMapApi) DeleteConfigMap(c *gin.Context) {
	req := configmap.DeleteConfigMapReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := k8sConfigMapService.DeleteConfigMap(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	} else {
		time.Sleep(1 * time.Second)
		response.OkWithMessage("删除成功", c)
	}
}

// CreateConfigMap 创建configmap
// @Tags kubernetes
// @Summary 创建configmap
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "创建configmap"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /kubernetes/configMap [post]
func (k *K8sConfigMapApi) CreateConfigMap(c *gin.Context) {
	req := configmap.CreateConfigMapReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if CronJob, err := k8sConfigMapService.CreateConfigMap(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(CronJob, c)
	}
}
