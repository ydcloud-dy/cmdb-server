package cicd

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cicd"
	"DYCLOUD/model/cicd/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type SourceCodeApi struct {
}

// GetSourceCodeList
//
//	@Description: 获取cicd 服务列表
//	@receiver s
//	@param c
func (s *SourceCodeApi) GetSourceCodeList(c *gin.Context) {
	var env *request.ServiceRequest
	err := c.BindQuery(&env)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//Environment.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, total, err := SourceCodeService.GetSourceCodeList(env)
	if err != nil {
		global.DYCLOUD_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     data,
		Total:    total,
		Page:     env.Page,
		PageSize: env.PageSize,
	}, "获取成功", c)
}

// CreateSourceCode
//
//	@Description: 创建服务列表
//	@receiver s
//	@param c
func (s *SourceCodeApi) CreateSourceCode(c *gin.Context) {
	var request *cicd.ServiceIntegration
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(request)
	request.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	err = SourceCodeService.CreateSourceCode(request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("创建成功", c)
}

// UpdateSourceCode
//
//	@Description: 更新服务
//	@receiver s
//	@param c
func (s *SourceCodeApi) UpdateSourceCode(c *gin.Context) {
	var request *cicd.ServiceIntegration
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(request, "================================")
	request.UpdatedBy = utils.GetUserID(c)
	err = SourceCodeService.UpdateSourceCode(request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("更新成功", c)
}

// DeleteSourceCode
//
//	@Description: 删除服务
//	@receiver s
//	@param c
func (s *SourceCodeApi) DeleteSourceCode(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(id)
	err = SourceCodeService.DeleteSourceCode(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("删除成功", c)
}
func (s *SourceCodeApi) DescribeSourceCode(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(id)
	//Environment.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, err := SourceCodeService.DescribeSourceCode(id)
	if err != nil {
		global.DYCLOUD_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

// VerifySourceCode
//
//	@Description: 验证服务是否正常连接
//	@receiver s
//	@param c
func (s *SourceCodeApi) VerifySourceCode(c *gin.Context) {
	var request *cicd.ServiceIntegration
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(request)
	request.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	msg, err := SourceCodeService.VerifySourceCode(request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(msg, c)
}
