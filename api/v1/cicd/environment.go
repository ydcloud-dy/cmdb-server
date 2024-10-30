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

type EnvironmentApi struct{}

// GetEnvironmentList
//
//	@Description: 获取环境列表
//	@receiver EnvironmentApi
//	@param c
func (EnvironmentApi *EnvironmentApi) GetEnvironmentList(c *gin.Context) {
	var env *request.EnvRequest
	err := c.BindQuery(&env)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(env)
	//Environment.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, total, err := EnvironmentService.GetEnvironmentList(env)
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

// DescribeEnvironment
//
//	@Description: 查看环境详情
//	@receiver EnvironmentApi
//	@param c
func (EnvironmentApi *EnvironmentApi) DescribeEnvironment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(id)
	//Environment.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, err := EnvironmentService.DescribeEnvironment(id)
	if err != nil {
		global.DYCLOUD_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

// CreateEnvironment
//
//	@Description: 创建环境
//	@receiver EnvironmentApi
//	@param c
func (EnvironmentApi *EnvironmentApi) CreateEnvironment(c *gin.Context) {
	var request *cicd.Environment
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(request)
	request.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	err = EnvironmentService.CreateEnvironment(request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("创建成功", c)
}

// UpdateEnvironment
//
//	@Description: 更新环境
//	@receiver EnvironmentApi
//	@param c
func (EnvironmentApi *EnvironmentApi) UpdateEnvironment(c *gin.Context) {
	var request *cicd.Environment
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(request, "================================")
	request.UpdatedBy = utils.GetUserID(c)
	data, err := EnvironmentService.UpdateEnvironment(request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

// DeleteEnvironment
//
//	@Description: 删除环境
//	@receiver EnvironmentApi
//	@param c
func (EnvironmentApi *EnvironmentApi) DeleteEnvironment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(id)
	err = EnvironmentService.DeleteEnvironment(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("删除成功", c)
}

// DeleteEnvironmentByIds
//
//	@Description: 根据id批量删除环境
//	@receiver EnvironmentApi
//	@param c
func (EnvironmentApi *EnvironmentApi) DeleteEnvironmentByIds(c *gin.Context) {
	ids := &request.DeleteEnvByIds{}
	err := c.ShouldBindJSON(ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(ids)
	err = EnvironmentService.DeleteEnvironmentByIds(ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("删除成功", c)
}
