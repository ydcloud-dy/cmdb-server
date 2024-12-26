package cmdb

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cmdb"
	cmdbReq "DYCLOUD/model/cmdb/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type CmdbHostsApi struct{}

// CreateCmdbHosts 创建cmdbHosts表
// @Tags CmdbHosts
// @Summary 创建cmdbHosts表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cmdb.CmdbHosts true "创建cmdbHosts表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /cmdbHosts/createCmdbHosts [post]
func (cmdbHostsApi *CmdbHostsApi) CreateCmdbHosts(c *gin.Context) {
	var cmdbHosts cmdb.CmdbHosts
	err := c.ShouldBindJSON(&cmdbHosts)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cmdbHosts.CreatedBy = utils.GetUserID(c)
	err = cmdbHostsService.CreateCmdbHosts(&cmdbHosts)
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// CreateCmdbHosts 创建cmdbHosts表
// @Tags CmdbHosts
// @Summary 创建cmdbHosts表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cmdb.CmdbHosts true "创建cmdbHosts表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /cmdbHosts/createCmdbHosts [post]
func (cmdbHostsApi *CmdbHostsApi) AuthenticationCmdbHosts(c *gin.Context) {
	var cmdbHosts cmdb.CmdbHosts
	err := c.ShouldBindJSON(&cmdbHosts)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cmdbHosts.CreatedBy = utils.GetUserID(c)
	err = cmdbHostsService.SSHTestCmdbHosts(&cmdbHosts)
	if err != nil {
		if err.Error() == "auth failed" {
			response.Result(177, nil, "auth failed", c)
			return
		}
		global.DYCLOUD_LOG.Error("验证失败!", zap.Error(err))
		response.FailWithMessage("验证失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("验证成功", c)
}

// DeleteCmdbHosts 删除cmdbHosts表
// @Tags CmdbHosts
// @Summary 删除cmdbHosts表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cmdb.CmdbHosts true "删除cmdbHosts表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /cmdbHosts/deleteCmdbHosts [delete]
func (cmdbHostsApi *CmdbHostsApi) DeleteCmdbHosts(c *gin.Context) {
	ID := c.Query("id")
	userID := utils.GetUserID(c)
	err := cmdbHostsService.DeleteCmdbHosts(ID, userID)
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteCmdbHostsByIds 批量删除cmdbHosts表
// @Tags CmdbHosts
// @Summary 批量删除cmdbHosts表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /cmdbHosts/deleteCmdbHostsByIds [delete]
func (cmdbHostsApi *CmdbHostsApi) DeleteCmdbHostsByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	err := cmdbHostsService.DeleteCmdbHostsByIds(IDs, userID)
	if err != nil {
		global.DYCLOUD_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateCmdbHosts 更新cmdbHosts表
// @Tags CmdbHosts
// @Summary 更新cmdbHosts表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cmdb.CmdbHosts true "更新cmdbHosts表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /cmdbHosts/updateCmdbHosts [put]
func (cmdbHostsApi *CmdbHostsApi) UpdateCmdbHosts(c *gin.Context) {
	var cmdbHosts cmdb.CmdbHosts
	err := c.ShouldBindJSON(&cmdbHosts)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cmdbHosts.UpdatedBy = utils.GetUserID(c)
	err = cmdbHostsService.UpdateCmdbHosts(cmdbHosts)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindCmdbHosts 用id查询cmdbHosts表
// @Tags CmdbHosts
// @Summary 用id查询cmdbHosts表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cmdb.CmdbHosts true "用id查询cmdbHosts表"
// @Success 200 {object} response.Response{data=cmdb.CmdbHosts,msg=string} "查询成功"
// @Router /cmdbHosts/findCmdbHosts [get]
func (cmdbHostsApi *CmdbHostsApi) FindCmdbHosts(c *gin.Context) {
	ID := c.Query("id")
	recmdbHosts, err := cmdbHostsService.GetCmdbHosts(ID)
	if err != nil {
		global.DYCLOUD_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(recmdbHosts, c)
}

// GetCmdbHostsList 分页获取cmdbHosts表列表
// @Tags CmdbHosts
// @Summary 分页获取cmdbHosts表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cmdbReq.CmdbHostsSearch true "分页获取cmdbHosts表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /cmdbHosts/getCmdbHostsList [get]
func (cmdbHostsApi *CmdbHostsApi) GetCmdbHostsList(c *gin.Context) {
	var pageInfo cmdbReq.CmdbHostsSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := cmdbHostsService.GetCmdbHostsInfoList(pageInfo)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// ImportHosts
//
//	@Description: 根据模板批量创建主机
//	@receiver cmdbHostsApi
//	@param c
func (cmdbHostsApi *CmdbHostsApi) ImportHosts(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file") // "file" 是前端上传文件的字段名
	if err != nil {
		response.FailWithMessage("获取文件失败: "+err.Error(), c)
		return
	}
	projectIdStr := c.PostForm("projectId")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		response.FailWithMessage("无效的 projectId: "+err.Error(), c)
		return
	}

	// 保存上传的文件到临时目录
	dst := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.FailWithMessage("保存文件失败: "+err.Error(), c)
		return
	}

	// 调用服务层处理上传逻辑
	if err := cmdbHostsService.ImportHosts(dst, projectId); err != nil {
		global.DYCLOUD_LOG.Error("导入失败!", zap.Error(err))
		response.FailWithMessage("导入失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("导入成功", c)
}

// GetCmdbHostsPublic 不需要鉴权的cmdbHosts表接口
// @Tags CmdbHosts
// @Summary 不需要鉴权的cmdbHosts表接口
// @accept application/json
// @Produce application/json
// @Param data query cmdbReq.CmdbHostsSearch true "分页获取cmdbHosts表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /cmdbHosts/getCmdbHostsPublic [get]
func (cmdbHostsApi *CmdbHostsApi) GetCmdbHostsPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	cmdbHostsService.GetCmdbHostsPublic()
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的cmdbHosts表接口信息",
	}, "获取成功", c)
}
