package configCenter

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/response"
	"DYCLOUD/model/configCenter"
	"DYCLOUD/model/configCenter/request"
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
//	@Description: 创建代码源
//	@receiver s
//	@param c
func (s *SourceCodeApi) CreateSourceCode(c *gin.Context) {
	var request *configCenter.ServiceIntegration
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
//	@Description: 更新代码源
//	@receiver s
//	@param c
func (s *SourceCodeApi) UpdateSourceCode(c *gin.Context) {
	var request *configCenter.ServiceIntegration
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
//	@Description: 删除代码源
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

// DescribeSourceCode
//
//	@Description: 查询代码源详细信息
//	@receiver s
//	@param c
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
//	@Description: 验证代码源是否正常连接
//	@receiver s
//	@param c
func (s *SourceCodeApi) VerifySourceCode(c *gin.Context) {
	var request *configCenter.ServiceIntegration
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

// GetGitProjectsByRepoId
//
//	@Description: 根据仓库id查询该仓库所有项目的列表
//	@receiver s
//	@param c
func (s *SourceCodeApi) GetGitProjectsByRepoId(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(id)
	//Environment.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, err := SourceCodeService.GetGitProjectsByRepoId(id)
	if err != nil {
		global.DYCLOUD_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}
