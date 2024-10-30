package api

import (
	"DYCLOUD/global"
	model "DYCLOUD/model/cloudCmdb"
	cloudcmdbreq "DYCLOUD/model/cloudCmdb/cloudcmdb"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/utils"
	cloudutils "DYCLOUD/utils/cloudCmdb"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CloudPlatformApi struct{}

// @Tags CloudPlatformApi
// @Summary 云厂商
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /cloudcmdb/cloud_platform/list [post]

func (p *CloudPlatformApi) CloudPlatformList(c *gin.Context) {
	var pageInfo cloudcmdbreq.SearchCloudPlatformParams
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	List, total, err := cloudPlatformService.List(pageInfo.CloudPlatform, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     List,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// @Tags CloudPlatformApi
// @Summary 根据id获取CloudPlatform
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "根据id获取api"
// @Success 200 {object} response.Response{data=systemRes.SysAPIResponse} "根据id获取api,返回包括api详情"
// @Router  /cloudcmdb/cloud_platform/getById [post]

func (p *CloudPlatformApi) GetCloudPlatformById(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cloud, regions, err := cloudPlatformService.GetCloudPlatformById(idInfo.ID)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(cloudcmdbreq.CloudResponse{CloudPlatform: cloud, Regions: regions}, "获取成功", c)
	}
}

// @Tags CloudPlatformApi
// @Summary 创建厂商
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.CloudPlatform true "创建厂商信息"
// @Success 200 {object} response.Response{msg=string} "创建厂商信息"
// @Router  /cloudcmdb/cloud_platform/create [post]

func (p *CloudPlatformApi) CreateCloudPlatform(c *gin.Context) {
	var cloud model.CloudPlatform
	_ = c.ShouldBindJSON(&cloud)
	if err := utils.Verify(cloud, cloudutils.CloudVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := cloudPlatformService.CreateCloudPlatform(cloud); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// @Tags CloudPlatformApi
// @Summary 更新厂商
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body  model.CloudPlatform true "更新"
// @Success 200 {object} response.Response{msg=string} "更新"
// @Router  /cloudcmdb/cloud_platform/update [put]

func (p *CloudPlatformApi) UpdateCloudPlatform(c *gin.Context) {
	var cloud model.CloudPlatform
	_ = c.ShouldBindJSON(&cloud)
	if err := cloudPlatformService.UpdateCloudPlatform(cloud); err != nil {
		global.DYCLOUD_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage("修改失败", c)
	} else {
		response.OkWithMessage("修改成功", c)
	}
}

// @Tags CloudPlatformApi
// @Summary 删除厂商
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.CloudPlatform true "ID"
// @Success 200 {object} response.Response{msg=string} "删除Clusters"
// @Router  /cloudcmdb/cloud_platform/delete [delete]

func (p *CloudPlatformApi) DeleteCloudPlatform(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := cloudPlatformService.DeleteCloudPlatform(idInfo); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags CloudPlatformApi
// @Summary 批量删除
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {object} response.Response{msg=string} "删除选中厂商"
// @Router  /cloudcmdb/cloud_platform/deleteByIds [delete]

func (p *CloudPlatformApi) DeleteCloudPlatformByIds(c *gin.Context) {
	var ids request.IdsReq
	_ = c.ShouldBindJSON(&ids)
	if err := cloudPlatformService.DeleteCloudPlatformByIds(ids); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}
