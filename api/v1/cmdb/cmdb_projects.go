package cmdb

import (
	
	"DYCLOUD/global"
    "DYCLOUD/model/common/response"
    "DYCLOUD/model/cmdb"
    cmdbReq "DYCLOUD/model/cmdb/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "DYCLOUD/utils"
)

type CmdbProjectsApi struct {}



// CreateCmdbProjects 创建cmdbProjects表
// @Tags CmdbProjects
// @Summary 创建cmdbProjects表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cmdb.CmdbProjects true "创建cmdbProjects表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /cmdbProjects/createCmdbProjects [post]
func (cmdbProjectsApi *CmdbProjectsApi) CreateCmdbProjects(c *gin.Context) {
	var cmdbProjects cmdb.CmdbProjects
	err := c.ShouldBindJSON(&cmdbProjects)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    cmdbProjects.CreatedBy = utils.GetUserID(c)
	err = cmdbProjectsService.CreateCmdbProjects(&cmdbProjects)
	if err != nil {
        global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteCmdbProjects 删除cmdbProjects表
// @Tags CmdbProjects
// @Summary 删除cmdbProjects表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cmdb.CmdbProjects true "删除cmdbProjects表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /cmdbProjects/deleteCmdbProjects [delete]
func (cmdbProjectsApi *CmdbProjectsApi) DeleteCmdbProjects(c *gin.Context) {
	ID := c.Query("ID")
    userID := utils.GetUserID(c)
	err := cmdbProjectsService.DeleteCmdbProjects(ID,userID)
	if err != nil {
        global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteCmdbProjectsByIds 批量删除cmdbProjects表
// @Tags CmdbProjects
// @Summary 批量删除cmdbProjects表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /cmdbProjects/deleteCmdbProjectsByIds [delete]
func (cmdbProjectsApi *CmdbProjectsApi) DeleteCmdbProjectsByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
    userID := utils.GetUserID(c)
	err := cmdbProjectsService.DeleteCmdbProjectsByIds(IDs,userID)
	if err != nil {
        global.DYCLOUD_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateCmdbProjects 更新cmdbProjects表
// @Tags CmdbProjects
// @Summary 更新cmdbProjects表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cmdb.CmdbProjects true "更新cmdbProjects表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /cmdbProjects/updateCmdbProjects [put]
func (cmdbProjectsApi *CmdbProjectsApi) UpdateCmdbProjects(c *gin.Context) {
	var cmdbProjects cmdb.CmdbProjects
	err := c.ShouldBindJSON(&cmdbProjects)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    cmdbProjects.UpdatedBy = utils.GetUserID(c)
	err = cmdbProjectsService.UpdateCmdbProjects(cmdbProjects)
	if err != nil {
        global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindCmdbProjects 用id查询cmdbProjects表
// @Tags CmdbProjects
// @Summary 用id查询cmdbProjects表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cmdb.CmdbProjects true "用id查询cmdbProjects表"
// @Success 200 {object} response.Response{data=cmdb.CmdbProjects,msg=string} "查询成功"
// @Router /cmdbProjects/findCmdbProjects [get]
func (cmdbProjectsApi *CmdbProjectsApi) FindCmdbProjects(c *gin.Context) {
	ID := c.Query("ID")
	recmdbProjects, err := cmdbProjectsService.GetCmdbProjects(ID)
	if err != nil {
        global.DYCLOUD_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(recmdbProjects, c)
}

// GetCmdbProjectsList 分页获取cmdbProjects表列表
// @Tags CmdbProjects
// @Summary 分页获取cmdbProjects表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cmdbReq.CmdbProjectsSearch true "分页获取cmdbProjects表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /cmdbProjects/getCmdbProjectsList [get]
func (cmdbProjectsApi *CmdbProjectsApi) GetCmdbProjectsList(c *gin.Context) {
	var pageInfo cmdbReq.CmdbProjectsSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := cmdbProjectsService.GetCmdbProjectsInfoList(pageInfo)
	if err != nil {
	    global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败:" + err.Error(), c)
        return
    }
    response.OkWithDetailed(response.PageResult{
        List:     list,
        Total:    total,
        Page:     pageInfo.Page,
        PageSize: pageInfo.PageSize,
    }, "获取成功", c)
}

// GetCmdbProjectsPublic 不需要鉴权的cmdbProjects表接口
// @Tags CmdbProjects
// @Summary 不需要鉴权的cmdbProjects表接口
// @accept application/json
// @Produce application/json
// @Param data query cmdbReq.CmdbProjectsSearch true "分页获取cmdbProjects表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /cmdbProjects/getCmdbProjectsPublic [get]
func (cmdbProjectsApi *CmdbProjectsApi) GetCmdbProjectsPublic(c *gin.Context) {
    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    cmdbProjectsService.GetCmdbProjectsPublic()
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的cmdbProjects表接口信息",
    }, "获取成功", c)
}
