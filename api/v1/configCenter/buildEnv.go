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

type BuildEnvApi struct{}

// GetBuildEnvList
//
//	@Description: 获取环境列表
//	@receiver BuildEnvApi
//	@param c
func (BuildEnvApi *BuildEnvApi) GetBuildEnvList(c *gin.Context) {
	var env request.BuildEnvRequest
	err := c.BindQuery(&env)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(env)
	//BuildEnv.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, total, err := BuildEnvService.GetBuildEnvList(&env)
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

// DescribeBuildEnv
//
//	@Description: 查看环境详情
//	@receiver BuildEnvApi
//	@param c
func (BuildEnvApi *BuildEnvApi) DescribeBuildEnv(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(id)
	//BuildEnv.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, err := BuildEnvService.DescribeBuildEnv(id)
	if err != nil {
		global.DYCLOUD_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

// CreateBuildEnv
//
//	@Description: 创建环境
//	@receiver BuildEnvApi
//	@param c
func (BuildEnvApi *BuildEnvApi) CreateBuildEnv(c *gin.Context) {
	var request *configCenter.BuildEnv
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(request)
	request.CreatedBy = utils.GetUserID(c)
	request.CreatedName = utils.GetUserName(c)
	//userId := utils.GetUserID(c)
	err = BuildEnvService.CreateBuildEnv(request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("创建成功", c)
}

// UpdateBuildEnv
//
//	@Description: 更新环境
//	@receiver BuildEnvApi
//	@param c
func (BuildEnvApi *BuildEnvApi) UpdateBuildEnv(c *gin.Context) {
	var request *configCenter.BuildEnv
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(request, "================================")
	request.UpdatedBy = utils.GetUserID(c)
	request.UpdatedName = utils.GetUserName(c)
	data, err := BuildEnvService.UpdateBuildEnv(request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

// DeleteBuildEnv
//
//	@Description: 删除环境
//	@receiver BuildEnvApi
//	@param c
func (BuildEnvApi *BuildEnvApi) DeleteBuildEnv(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(id)
	err = BuildEnvService.DeleteBuildEnv(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("删除成功", c)
}

// DeleteBuildEnvByIds
//
//	@Description: 根据id批量删除环境
//	@receiver BuildEnvApi
//	@param c
func (BuildEnvApi *BuildEnvApi) DeleteBuildEnvByIds(c *gin.Context) {
	ids := &request.DeleteBuildEnvByIds{}
	err := c.ShouldBindJSON(ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(ids)
	err = BuildEnvService.DeleteBuildEnvByIds(ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("删除成功", c)
}
