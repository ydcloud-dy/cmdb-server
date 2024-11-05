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

type ServiceIntegrationApi struct {
}

// GetServiceIntegrationList
//
//	@Description: 获取cicd 服务列表
//	@receiver s
//	@param c
func (s *ServiceIntegrationApi) GetServiceIntegrationList(c *gin.Context) {
	var env *request.ServiceRequest
	err := c.BindQuery(&env)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//Environment.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, total, err := ServiceIntegrationService.GetServiceIntegrationList(env)
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

// CreateServiceIntegration
//
//	@Description: 创建服务列表
//	@receiver s
//	@param c
func (s *ServiceIntegrationApi) CreateServiceIntegration(c *gin.Context) {
	var request *configCenter.ServiceIntegration
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(request)
	request.CreatedBy = utils.GetUserID(c)
	request.CreatedName = utils.GetUserName(c)
	//userId := utils.GetUserID(c)
	err = ServiceIntegrationService.CreateServiceIntegration(request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData("创建成功", c)
}

// UpdateServiceIntegration
//
//	@Description: 更新服务
//	@receiver s
//	@param c
func (s *ServiceIntegrationApi) UpdateServiceIntegration(c *gin.Context) {
	var request *configCenter.ServiceIntegration
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(request, "================================")
	request.UpdatedBy = utils.GetUserID(c)
	request.UpdatedName = utils.GetUserName(c)
	err = ServiceIntegrationService.UpdateServiceIntegration(request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("更新成功", c)
}

// DeleteServiceIntegration
//
//	@Description: 删除服务
//	@receiver s
//	@param c
func (s *ServiceIntegrationApi) DeleteServiceIntegration(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(id)
	err = ServiceIntegrationService.DeleteServiceIntegration(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("删除成功", c)
}

// DescribeServiceIntegration
//
//	@Description: 查询服务详细信息
//	@receiver s
//	@param c
func (s *ServiceIntegrationApi) DescribeServiceIntegration(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(id)
	//Environment.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, err := ServiceIntegrationService.DescribeServiceIntegration(id)
	if err != nil {
		global.DYCLOUD_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

// VerifyServiceIntegration
//
//	@Description: 验证服务是否正常连接
//	@receiver s
//	@param c
func (s *ServiceIntegrationApi) VerifyServiceIntegration(c *gin.Context) {
	var request *configCenter.ServiceIntegration
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(request)
	request.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	msg, err := ServiceIntegrationService.VerifyServiceIntegration(request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(msg, c)
}
