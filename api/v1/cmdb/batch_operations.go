package cmdb

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cmdb/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BatchOperationsApi struct{}

func (BatchOperationsApi *BatchOperationsApi) ExecuteCommands(c *gin.Context) {
	var BatchOperations request.ExecuteRequest
	err := c.ShouldBindJSON(&BatchOperations)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//BatchOperations.CreatedBy = utils.GetUserID(c)
	userId := utils.GetUserID(c)
	BatchOperations.UserId = userId
	data, err := batchOperationsService.CreateBatchOperations(BatchOperations)
	if err != nil {
		global.DYCLOUD_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}
func (BatchOperationsApi *BatchOperationsApi) ExecuteRecords(c *gin.Context) {
	userId := utils.GetUserID(c)
	// 调用 service 获取当前用户最近的10条记录
	data, err := batchOperationsService.GetUserRecentExecutionRecords(userId)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取执行记录失败!", zap.Error(err))
		response.FailWithMessage("获取执行记录失败:"+err.Error(), c)
		return
	}

	// 成功返回数据
	response.OkWithData(data, c)
}
