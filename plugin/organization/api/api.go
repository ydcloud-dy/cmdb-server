package api

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	organization "DYCLOUD/plugin/organization/model"
	organizationReq "DYCLOUD/plugin/organization/model/request"
	"DYCLOUD/plugin/organization/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrganizationApi struct {
}

var orgService = service.ServiceGroupApp.OrganizationService

// CreateOrganization 创建Organization
// @Tags Organization
// @Summary 创建Organization
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body organization.Organization true "创建Organization"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /org/createOrganization [post]
func (orgApi *OrganizationApi) CreateOrganization(c *gin.Context) {
	var org organization.Organization
	_ = c.ShouldBindJSON(&org)
	if err := orgService.CreateOrganization(org); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteOrganization 删除Organization
// @Tags Organization
// @Summary 删除Organization
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body organization.Organization true "删除Organization"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /org/deleteOrganization [delete]
func (orgApi *OrganizationApi) DeleteOrganization(c *gin.Context) {
	var org organization.Organization
	_ = c.ShouldBindJSON(&org)
	if err := orgService.DeleteOrganization(org); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteOrganizationByIds 批量删除Organization
// @Tags Organization
// @Summary 批量删除Organization
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Organization"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /org/deleteOrganizationByIds [delete]
func (orgApi *OrganizationApi) DeleteOrganizationByIds(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := orgService.DeleteOrganizationByIds(IDS); err != nil {
		global.DYCLOUD_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateOrganization 更新Organization
// @Tags Organization
// @Summary 更新Organization
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body organization.Organization true "更新Organization"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /org/updateOrganization [put]
func (orgApi *OrganizationApi) UpdateOrganization(c *gin.Context) {
	var org organization.Organization
	_ = c.ShouldBindJSON(&org)
	if err := orgService.UpdateOrganization(org); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindOrganization 用id查询Organization
// @Tags Organization
// @Summary 用id查询Organization
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query organization.Organization true "用id查询Organization"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /org/findOrganization [get]
func (orgApi *OrganizationApi) FindOrganization(c *gin.Context) {
	var org organization.Organization
	_ = c.ShouldBindQuery(&org)
	if reorg, err := orgService.GetOrganization(org.ID); err != nil {
		global.DYCLOUD_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reorg": reorg}, c)
	}
}

// GetOrganizationList 分页获取Organization列表
// @Tags Organization
// @Summary 分页获取Organization列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query organizationReq.OrganizationSearch true "分页获取Organization列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /org/getOrganizationList [get]
func (orgApi *OrganizationApi) GetOrganizationList(c *gin.Context) {
	var pageInfo organizationReq.OrganizationSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if list, total, err := orgService.GetOrganizationInfoList(pageInfo); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetOrganizationList 获取Organization树
// @Tags Organization
// @Summary 获取Organization树
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query organizationReq.OrganizationSearch true "获取Organization树"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /org/getOrganizationTree [get]
func (orgApi *OrganizationApi) GetOrganizationTree(c *gin.Context) {
	if tree, err := orgService.GetOrganizationTree(); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(tree, "获取成功", c)
	}
}

func (orgApi *OrganizationApi) CreateOrgUser(c *gin.Context) {
	var orgUser organization.OrgUserReq
	c.ShouldBindJSON(&orgUser)
	if err := orgService.CreateOrgUser(orgUser); err != nil {
		global.DYCLOUD_LOG.Error("变更失败!", zap.Error(err))
		response.FailWithMessage("变更失败", c)
	} else {
		response.OkWithMessage("变更成功", c)
	}
}

func (orgApi *OrganizationApi) FindOrgUserAll(c *gin.Context) {
	org := c.Query("organizationID")
	if UserIds, err := orgService.FindOrgUserAll(org); err != nil {
		global.DYCLOUD_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(UserIds, c)
	}
}

func (orgApi *OrganizationApi) FindOrgUserList(c *gin.Context) {
	var pageInfo organizationReq.OrgUserSearch
	c.ShouldBindQuery(&pageInfo)
	if list, total, err := orgService.GetOrgUserList(pageInfo); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

func (orgApi *OrganizationApi) SetOrgUserAdmin(c *gin.Context) {
	var orgUser organization.OrgUser
	c.ShouldBindJSON(&orgUser)
	if err := orgService.SetOrgUserAdmin(orgUser.SysUserID, orgUser.IsAdmin); err != nil {
		global.DYCLOUD_LOG.Error("设置失败!", zap.Error(err))
		response.FailWithMessage("设置失败", c)
	} else {
		response.OkWithMessage("设置成功", c)
	}
}

func (orgApi *OrganizationApi) SetOrgAuthority(c *gin.Context) {
	var dataAuthority organization.DataAuthority
	c.ShouldBindJSON(&dataAuthority)
	if err := orgService.SetOrgAuthority(dataAuthority.AuthorityID, dataAuthority.AuthorityType); err != nil {
		global.DYCLOUD_LOG.Error("设置失败!", zap.Error(err))
		response.FailWithMessage("设置失败", c)
	} else {
		response.OkWithMessage("设置成功", c)
	}
}

func (orgApi *OrganizationApi) SyncAuthority(c *gin.Context) {
	if err := orgService.SyncAuthority(); err != nil {
		global.DYCLOUD_LOG.Error("同步失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("同步成功", c)
	}
}

func (orgApi *OrganizationApi) GetAuthority(c *gin.Context) {
	if authDataList, err := orgService.GetOrgAuthority(); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithData(authDataList, c)
	}
}

func (orgApi *OrganizationApi) DeleteOrgUser(c *gin.Context) {
	var orgUser organization.OrgUserReq
	c.ShouldBindJSON(&orgUser)
	if err := orgService.DeleteOrgUser(orgUser.SysUserIDS, orgUser.OrganizationID); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

func (orgApi *OrganizationApi) TransferOrgUser(c *gin.Context) {
	var orgUser organization.OrgUserReq
	c.ShouldBindJSON(&orgUser)
	if err := orgService.TransferOrgUser(orgUser.SysUserIDS, orgUser.OrganizationID, orgUser.ToOrganizationID); err != nil {
		global.DYCLOUD_LOG.Error("转移失败!", zap.Error(err))
		response.FailWithMessage("转移失败", c)
	} else {
		response.OkWithMessage("转移成功", c)
	}
}
