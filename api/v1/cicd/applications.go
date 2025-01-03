package cicd

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cicd"
	request2 "DYCLOUD/model/cicd/request"
	request3 "DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/model/kubernetes/pods"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type ApplicationsApi struct{}

// GetApplicationsList
//
//	@Description: 获取应用列表
//	@receiver ApplicationsApi
//	@param c
func (ApplicationsApi *ApplicationsApi) GetApplicationsList(c *gin.Context) {
	var env *request2.ApplicationRequest
	err := c.BindQuery(&env)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//Applications.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, total, err := ApplicationService.GetApplicationsList(env)
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
func (ApplicationsApi *ApplicationsApi) GetBranchesList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var branches = request2.ApplicationRequest{}
	err = c.BindQuery(&branches)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	branches.AppId = id

	//Applications.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, total, err := ApplicationService.GetBranchesList(branches)
	if err != nil {
		global.DYCLOUD_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     data,
		Total:    total,
		Page:     branches.Page,
		PageSize: branches.PageSize,
	}, "获取成功", c)
}

// DescribeApplications
//
//	@Description: 查看环境详情
//	@receiver ApplicationsApi
//	@param c
func (ApplicationsApi *ApplicationsApi) DescribeApplications(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//Applications.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, err := ApplicationService.DescribeApplications(id)
	if err != nil {
		global.DYCLOUD_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}
func (ApplicationsApi *ApplicationsApi) DescribeApplicationsByName(c *gin.Context) {
	name := c.Query("app_code")

	//Applications.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, err := ApplicationService.DescribeApplicationsByName(name)
	if err != nil {
		global.DYCLOUD_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

// CreateApplications
//
//	@Description: 创建环境
//	@receiver ApplicationsApi
//	@param c
func (ApplicationsApi *ApplicationsApi) CreateApplications(c *gin.Context) {

	var request *cicd.AppRequestBody
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	request.App.CreatedBy = utils.GetUserID(c)
	request.App.CreatedName = utils.GetUserName(c)
	//userId := utils.GetUserID(c)
	err = ApplicationService.CreateApplications(request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("创建成功", c)
}
func (ApplicationsApi *ApplicationsApi) SyncBranches(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//Applications.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	err = ApplicationService.SyncBranches(id)
	if err != nil {
		global.DYCLOUD_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithData("同步成功", c)
}

// UpdateApplications
//
//	@Description: 更新环境
//	@receiver ApplicationsApi
//	@param c
func (ApplicationsApi *ApplicationsApi) UpdateApplications(c *gin.Context) {
	var request *cicd.AppRequestBody
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	request.App.UpdatedBy = utils.GetUserID(c)
	request.App.UpdatedName = utils.GetUserName(c)
	data, err := ApplicationService.UpdateApplications(request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

// DeleteApplications
//
//	@Description: 删除环境
//	@receiver ApplicationsApi
//	@param c
func (ApplicationsApi *ApplicationsApi) DeleteApplications(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = ApplicationService.DeleteApplications(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("删除成功", c)
}

// DeleteApplicationsByIds
//
//	@Description: 根据id批量删除环境
//	@receiver ApplicationsApi
//	@param c
func (ApplicationsApi *ApplicationsApi) DeleteApplicationsByIds(c *gin.Context) {
	ids := &request2.DeleteApplicationByIds{}
	err := c.ShouldBindJSON(ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = ApplicationService.DeleteApplicationsByIds(ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("删除成功", c)
}

func (ApplicationsApi *ApplicationsApi) GetApplicationDeploymentInfo(c *gin.Context) {
	request := request2.DeploymentInfoRequest{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	podlist, total, err := ApplicationService.GetApplicationDeploymentInfo(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(pods.PodListResponse{
		Items: podlist,
		Total: total,
		PageInfo: request3.PageInfo{
			Page:     1,
			PageSize: 10000,
		},
	}, "获取成功", c)
}
