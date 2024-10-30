package api

import (
	"DYCLOUD/global"
	cloudcmdbreq "DYCLOUD/model/cloudCmdb/cloudcmdb"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CloudRDSApi struct{}

// @Tags CloudRDSApi
// @Summary 同步RDS
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /cloudcmdb/rds/sync [post]

func (r *CloudRDSApi) CloudRDSSync(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := cloudRDSService.SyncRDS(idInfo.ID)
	if err != nil {
		global.DYCLOUD_LOG.Error("同步操作失败!", zap.Error(err))
		response.FailWithMessage("同步操作失败", c)
	} else {
		response.OkWithMessage("同步操作成功, 数据异步处理中, 请稍后!", c)
	}
}

// @Tags CloudRDSApi
// @Summary RDS列表
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /cloudcmdb/rds/list [post]

func (r *CloudRDSApi) CloudRDSList(c *gin.Context) {
	var pageInfo cloudcmdbreq.SearchRDSParams
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	List, total, err := cloudRDSService.List(pageInfo.RDS, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
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

// @Tags CloudRDSApi
// @Summary RDs Tree数据
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /cloudcmdb/rds/tree [post]

func (r *CloudRDSApi) CloudRDSTree(c *gin.Context) {
	var pageInfo cloudcmdbreq.SearchCloudPlatformParams
	_ = c.ShouldBindJSON(&pageInfo)
	list, err := cloudRDSService.RDSTree(pageInfo.CloudPlatform, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取目录树失败!", zap.Error(err))
		response.FailWithMessage("获取目录树失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List: list,
		}, "获取成功", c)
	}
}
