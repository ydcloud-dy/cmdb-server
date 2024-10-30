package velero

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	veleroReq "DYCLOUD/model/velero/request"
	"DYCLOUD/service"
	velero2 "DYCLOUD/service/kubernetes/velero"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type K8sVeleroTasksApi struct{}

var k8sVeleroService = service.ServiceGroupApp.VeleroServiceGroup.K8sVeleroService

// CreateK8sVeleroTasks 创建 velero 备份任务
// @Summary 创建 velero 备份任务
// @Description 创建 velero 备份任务
// @Tags Velero
// @Accept  json
// @Produce  json
// @Param data body veleroReq.CreateVeleroTaskReq true "创建 velero 备份任务"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "创建成功"
// @Router /velero/tasks [post]
func (k8sVeleroTasksApi *K8sVeleroTasksApi) CreateK8sVeleroTasks(c *gin.Context) {
	var k8sVeleroTasks veleroReq.CreateVeleroTaskReq
	err := c.ShouldBindJSON(&k8sVeleroTasks)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	service := velero2.K8sVeleroTasksService{}
	ins, err := service.CreateK8sVeleroTasks(&k8sVeleroTasks, utils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(ins, c)
}

// DeleteK8sVeleroTasks 删除 velero 备份任务
// @Summary 删除 velero 备份任务
// @Description 删除 velero 备份任务
// @Tags Velero
// @Accept  json
// @Produce  json
// @Param data body veleroReq.DeleteVeleroTaskReq true "删除 velero 备份任务"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /velero/tasks [delete]
func (k8sVeleroTasksApi *K8sVeleroTasksApi) DeleteK8sVeleroTasks(c *gin.Context) {
	var req veleroReq.DeleteVeleroTaskReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	service := velero2.K8sVeleroTasksService{}
	err = service.DeleteK8sVeleroTasks(&req, utils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateK8sVeleroTasks 更新 velero 任务
// @Summary 更新 velero 任务
// @Description 更新 velero 任务
// @Tags Velero
// @Accept  json
// @Produce  json
// @Param data body veleroReq.UpdateVeleroTaskReq true "更新 velero 任务"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "更新成功"
// @Router /velero/tasks [put]
func (k8sVeleroTasksApi *K8sVeleroTasksApi) UpdateK8sVeleroTasks(c *gin.Context) {
	req := veleroReq.UpdateVeleroTaskReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	service := velero2.K8sVeleroTasksService{}

	if list, err := service.UpdateK8sVeleroTasks(&req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}

// DescribeVeleroTasks 查看 velero 任务详情
// @Summary 查看 velero 任务详情
// @Description 查看 velero 任务详情
// @Tags Velero
// @Accept  json
// @Produce  json
// @Param data query veleroReq.DescribeVeleroTaskReq true "查看 velero 任务详情"
// @Success 200 {object} response.Response{data=veleroReq.DescribeK8sVeleroTaskResponse,msg=string} "获取成功"
// @Router /velero/taskDetail [get]
func (k8sVeleroTasksApi *K8sVeleroTasksApi) DescribeVeleroTasks(c *gin.Context) {
	req := veleroReq.DescribeVeleroTaskReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	service := velero2.K8sVeleroTasksService{}

	veleroTask, err := service.DescribeVeleroTask(&req, utils.GetUserUuid(c))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(veleroReq.DescribeK8sVeleroTaskResponse{Items: veleroTask}, "获取成功", c)
}

// GetK8sVeleroTasksList 获取 velero 任务列表
// @Summary 获取 velero 任务列表
// @Description 获取 velero 任务列表
// @Tags Velero
// @Accept  json
// @Produce  json
// @Param data query veleroReq.K8sVeleroTasksSearchReq true "获取 velero 任务列表"
// @Success 200 {object} response.Response{data=veleroReq.K8sVeleroTaskListResponse,msg=string} "获取成功"
// @Router /velero/tasks [get]
func (k8sVeleroTasksApi *K8sVeleroTasksApi) GetK8sVeleroTasksList(c *gin.Context) {
	var pageInfo veleroReq.K8sVeleroTasksSearchReq
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	service := velero2.K8sVeleroTasksService{}

	list, total, err := service.GetK8sVeleroTasksInfoList(pageInfo, utils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(veleroReq.K8sVeleroTaskListResponse{
		Items: list,
		Total: total,
		PageInfo: request.PageInfo{
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
			Keyword:  pageInfo.Keyword,
		},
	}, "获取成功", c)
}

// CreateVelero 创建 velero
// @Summary 创建 velero
// @Description 创建 velero
// @Tags Velero
// @Accept  json
// @Produce  json
// @Param data body veleroReq.VeleroModel true "创建 velero"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /velero [post]
func (k8sVeleroTasksApi *K8sVeleroTasksApi) CreateVelero(c *gin.Context) {
	req := veleroReq.VeleroModel{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := k8sVeleroService.CreateVelero(&req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithMessage("velero创建成功！", c)
	}
}

// ReductionK8sVeleroRecord 恢复 velero 记录
// @Summary 恢复 velero 记录
// @Description 恢复 velero 记录
// @Tags Velero
// @Accept  json
// @Produce  json
// @Param data body veleroReq.CreateVeleroRestoreReq true "恢复 velero 记录"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "创建成功"
// @Router /velero/record/reduction [post]
func (K8sVeleroTasksApi *K8sVeleroTasksApi) ReductionK8sVeleroRecord(c *gin.Context) {
	var k8sVeleroTasks veleroReq.CreateVeleroRestoreReq
	err := c.ShouldBindJSON(&k8sVeleroTasks)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	service := velero2.K8sVeleroRestoresService{}
	ins, err := service.CreateK8sVeleroRestore(&k8sVeleroTasks, utils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(ins, c)
}

// GetK8sVeleroRecordList 获取 velero 备份记录列表
// @Summary 获取 velero 备份记录列表
// @Description 获取 velero 备份记录列表
// @Tags Velero
// @Accept  json
// @Produce  json
// @Param data query veleroReq.K8sVeleroRecordsSearchReq true "获取 velero 备份记录列表"
// @Success 200 {object} response.Response{data=veleroReq.K8sVeleroRecordListResponse,msg=string} "获取成功"
// @Router /velero/record [get]
func (K8sVeleroTasksApi *K8sVeleroTasksApi) GetK8sVeleroRecordList(c *gin.Context) {
	var req veleroReq.K8sVeleroRecordsSearchReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	service := velero2.K8sVeleroRecordsService{}

	list, total, err := service.GetK8sVeleroRecordList(req, utils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(veleroReq.K8sVeleroRecordListResponse{
		Items: list,
		Total: total,
		PageInfo: request.PageInfo{
			Page:     req.Page,
			PageSize: req.PageSize,
			Keyword:  req.Keyword,
		},
	}, "获取成功", c)
}

// DeleteK8sVeleroRecord 删除 velero 备份记录
// @Summary 删除 velero 备份记录
// @Description 删除 velero 备份记录
// @Tags Velero
// @Accept  json
// @Produce  json
// @Param data body veleroReq.DeleteVeleroRecordReq true "删除 velero 备份记录"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /velero/record [delete]
func (K8sVeleroTasksApi *K8sVeleroTasksApi) DeleteK8sVeleroRecord(c *gin.Context) {
	var req veleroReq.DeleteVeleroRecordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	service := velero2.K8sVeleroRecordsService{}
	err = service.DeleteK8sVeleroRecord(&req, utils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DescribeK8sVeleroRecord 查看 velero 备份详情
// @Summary 查看 velero 备份详情
// @Description 查看 velero 备份详情
// @Tags Velero
// @Accept  json
// @Produce  json
// @Param data query veleroReq.DescribeVeleroRecordReq true "查看 velero 备份详情"
// @Success 200 {object} response.Response{data=veleroReq.DescribeK8sVeleroRecordResponse,msg=string} "获取成功"
// @Router /velero/recordDetail [get]
func (K8sVeleroTasksApi *K8sVeleroTasksApi) DescribeK8sVeleroRecord(c *gin.Context) {
	req := veleroReq.DescribeVeleroRecordReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	service := velero2.K8sVeleroRecordsService{}

	veleroBackup, err := service.DescribeVeleroRecord(&req, utils.GetUserUuid(c))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(veleroReq.DescribeK8sVeleroRecordResponse{Items: veleroBackup}, "获取成功", c)
}

// CreateK8sVeleroRecord 创建 velero 备份
// @Summary 创建 velero 备份
// @Description 创建 velero 备份
// @Tags Velero
// @Accept  json
// @Produce  json
// @Param data body veleroReq.CreateVeleroRecordReq true "创建 velero 备份"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "创建成功"
// @Router /velero/record [post]
func (k8sVeleroTasksApi *K8sVeleroTasksApi) CreateK8sVeleroRecord(c *gin.Context) {
	var k8sVeleroTasks veleroReq.CreateVeleroRecordReq
	err := c.ShouldBindJSON(&k8sVeleroTasks)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	service := velero2.K8sVeleroRecordsService{}
	ins, err := service.CreateK8sVeleroRecord(&k8sVeleroTasks, utils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(ins, c)
}

// UpdateK8sVeleroRecord 更新 velero 备份记录
// @Summary 更新 velero 备份记录
// @Description 更新 velero 备份记录
// @Tags Velero
// @Accept  json
// @Produce  json
// @Param data body veleroReq.UpdateVeleroRecordReq true "更新 velero 备份记录"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "更新成功"
// @Router /velero/record [put]
func (k8sVeleroTasksApi *K8sVeleroTasksApi) UpdateK8sVeleroRecord(c *gin.Context) {
	req := veleroReq.UpdateVeleroRecordReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	service := velero2.K8sVeleroRecordsService{}

	if list, err := service.UpdateK8sVeleroRecord(&req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}

// GetK8sVeleroRestoreList 获取 velero 回滚记录列表
// @Summary 获取 velero 回滚记录列表
// @Description 获取 velero 回滚记录列表
// @Tags Velero
// @Accept  json
// @Produce  json
// @Param data query veleroReq.K8sVeleroRestoresSearchReq true "获取 velero 回滚记录列表"
// @Success 200 {object} response.Response{data=veleroReq.K8sVeleroRestoreListResponse,msg=string} "获取成功"
// @Router /velero/restore [get]
func (K8sVeleroTasksApi *K8sVeleroTasksApi) GetK8sVeleroRestoreList(c *gin.Context) {
	var req veleroReq.K8sVeleroRestoresSearchReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	service := velero2.K8sVeleroRestoresService{}

	list, total, err := service.GetK8sVeleroRestoreList(req, utils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(veleroReq.K8sVeleroRestoreListResponse{
		Items: list,
		Total: total,
		PageInfo: request.PageInfo{
			Page:     req.Page,
			PageSize: req.PageSize,
			Keyword:  req.Keyword,
		},
	}, "获取成功", c)
}

// DeleteK8sVeleroRestore 删除 velero 回滚记录
// @Summary 删除 velero 回滚记录
// @Description 删除 velero 回滚记录
// @Tags Velero
// @Accept  json
// @Produce  json
// @Param data body veleroReq.DeleteVeleroRestoreReq true "删除 velero 回滚记录"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /velero/restore [delete]
func (K8sVeleroTasksApi *K8sVeleroTasksApi) DeleteK8sVeleroRestore(c *gin.Context) {
	var req veleroReq.DeleteVeleroRestoreReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	service := velero2.K8sVeleroRestoresService{}
	err = service.DeleteK8sVeleroRestore(&req, utils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	time.Sleep(time.Second * 1)
	response.OkWithMessage("删除成功", c)
}

// DescribeK8sVeleroRestore 查看 velero 回滚记录详情
// @Summary 查看 velero 回滚记录详情
// @Description 查看 velero 回滚记录详情
// @Tags Velero
// @Accept  json
// @Produce  json
// @Param data query veleroReq.DescribeVeleroRestoreReq true "查看 velero 回滚记录详情"
// @Success 200 {object} response.Response{data=veleroReq.DescribeK8sVeleroRestoreResponse,msg=string} "获取成功"
// @Router /velero/restoreDetail [get]
func (K8sVeleroTasksApi *K8sVeleroTasksApi) DescribeK8sVeleroRestore(c *gin.Context) {
	req := veleroReq.DescribeVeleroRestoreReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	service := velero2.K8sVeleroRestoresService{}

	veleroRestore, err := service.DescribeVeleroRestore(&req, utils.GetUserUuid(c))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(veleroReq.DescribeK8sVeleroRestoreResponse{Items: veleroRestore}, "获取成功", c)
}
