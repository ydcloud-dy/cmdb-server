package api

import (
	config "DYCLOUD/config"
	"DYCLOUD/global"
	"DYCLOUD/model/common/response"
	model "DYCLOUD/model/docker"
	service "DYCLOUD/service/docker"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HostApi struct {
}

var hsService = service.ServiceGroupApp.HostService

// CreateHost 创建主机列表
// @Tags Host
// @Summary 创建主机列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body host.Host true "创建主机列表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /hs/createHost [post]
func (hsApi *HostApi) CreateHost(c *gin.Context) {
	var hs model.Host
	err := c.ShouldBindJSON(&hs)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	//if strings.Trim(hs.Name, " ") == "host" {
	//	response.FailWithMessage("主机名称不能使用host宿主机的名称", c)
	//	return
	//}

	hs.CreatedBy = utils.GetUserID(c)

	if err := hsService.CreateHost(&hs); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteHost 删除主机列表
// @Tags Host
// @Summary 删除主机列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body host.Host true "删除主机列表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /hs/deleteHost [delete]
func (hsApi *HostApi) DeleteHost(c *gin.Context) {

	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	if err := hsService.DeleteHost(ID, userID); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteHostByIds 批量删除主机列表
// @Tags Host
// @Summary 批量删除主机列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /hs/deleteHostByIds [delete]
func (hsApi *HostApi) DeleteHostByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	if err := hsService.DeleteHostByIds(IDs, userID); err != nil {
		global.DYCLOUD_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateHost 更新主机列表
// @Tags Host
// @Summary 更新主机列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body host.Host true "更新主机列表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /hs/updateHost [put]
func (hsApi *HostApi) UpdateHost(c *gin.Context) {
	var hs model.Host
	err := c.ShouldBindJSON(&hs)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//if strings.Trim(hs.Name, " ") == "host" {
	//	response.FailWithMessage("不能修改host宿主机", c)
	//	return
	//}

	hs.UpdatedBy = utils.GetUserID(c)

	if err := hsService.UpdateHost(hs); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindHost 用id查询主机列表
// @Tags Host
// @Summary 用id查询主机列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query host.Host true "用id查询主机列表"
// @Success 200 {object} response.Response{data=object{rehs=host.Host},msg=string} "查询成功"
// @Router /hs/findHost [get]
func (hsApi *HostApi) FindHost(c *gin.Context) {
	ID := c.Query("ID")
	if rehs, err := hsService.GetHost(ID); err != nil {
		global.DYCLOUD_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rehs": rehs}, c)
	}
}

// GetHostList 分页获取主机列表列表
// @Tags Host
// @Summary 分页获取主机列表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query hostReq.HostSearch true "分页获取主机列表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /hs/getHostList [get]
func (hsApi *HostApi) GetHostList(c *gin.Context) {
	var pageInfo model.HostSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := hsService.GetHostInfoList(pageInfo); err != nil {
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

// GetHostPublic 不需要鉴权的主机列表接口
// @Tags Host
// @Summary 不需要鉴权的主机列表接口
// @accept application/json
// @Produce application/json
// @Param data query hostReq.HostSearch true "分页获取主机列表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /hs/getHostPublic [get]
func (hsApi *HostApi) GetHostPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的主机列表接口信息",
	}, "获取成功", c)
}

// GetHostGenTlsScript 获取主机生成ssl脚本
// @Tags Host
// @accept application/json
// @Produce application/json
// @Router /hs/getHostPublic [get]
func (hsApi *HostApi) GetHostGenTlsScript(c *gin.Context) {
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=gen-tls.sh") // 用来指定下载下来的文件名
	c.Header("Content-Transfer-Encoding", "binary")
	c.Writer.Write([]byte(config.GenHostScript))
}

// CheckHost 检测主机可用性
// @Tags Host
// @Summary 检测主机可用性
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body host.Host true "检测主机可用性"
// @Success 200 {object} response.Response{msg=string} "检测主机可用性"
// @Router /hs/updateHost [post]
func (hsApi *HostApi) CheckHost(c *gin.Context) {
	var hs model.Host
	err := c.ShouldBindJSON(&hs)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if hsRes, err := hsService.CheckHost(hs); err != nil {
		global.DYCLOUD_LOG.Error("检查docker连接可用性失败!", zap.Error(err))
		response.FailWithMessage("检查docker连接可用性失败", c)
	} else {
		response.OkWithData(hsRes, c)
	}
}
