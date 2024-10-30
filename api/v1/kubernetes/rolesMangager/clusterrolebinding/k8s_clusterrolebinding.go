package clusterrolebinding

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/model/kubernetes/clusterolebinding"
	"DYCLOUD/service"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type K8sClusterRoleBindingApi struct{}

var k8sClusterRoleBindingService = service.ServiceGroupApp.ClusterRoleBindingServiceGroup.K8sClusterRoleBindingService

// GetClusterRoleBindingList 获取集群角色绑定列表
// @Tags kubernetes
// @Summary 获取集群角色绑定列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "获取集群角色绑定列表"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /kubernetes/ClusterRoleBinding [get]
func (k *K8sClusterRoleBindingApi) GetClusterRoleBindingList(c *gin.Context) {
	req := clusterolebinding.GetClusterRoleBindingListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := k8sClusterRoleBindingService.GetClusterRoleBindingList(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	} else {

		response.OkWithDetailed(clusterolebinding.ClusterRoleBindingListResponse{
			Items: list,
			Total: total,
			PageInfo: request.PageInfo{
				Page:     req.Page,
				PageSize: req.PageSize,
				Keyword:  req.Keyword,
			},
		}, "获取成功", c)
	}
}

// DescribeClusterRoleBindingInfo 获取集群角色绑定详情
// @Tags kubernetes
// @Summary 获取集群角色绑定详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "获取集群角色绑定详情"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /kubernetes/ClusterRoleBindingDetails [get]
func (k *K8sClusterRoleBindingApi) DescribeClusterRoleBindingInfo(c *gin.Context) {
	req := clusterolebinding.DescribeClusterRoleBindingReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sClusterRoleBindingService.DescribeClusterRoleBinding(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(clusterolebinding.DescribeClusterRoleBindingResponse{Items: list}, "获取成功", c)
	}
}

// UpdateClusterRoleBinding 更新集群角色绑定详情
// @Tags kubernetes
// @Summary 更新集群角色绑定详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "更新集群角色绑定详情"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /kubernetes/ClusterRoleBinding [put]
func (k *K8sClusterRoleBindingApi) UpdateClusterRoleBinding(c *gin.Context) {
	req := clusterolebinding.UpdateClusterRoleBindingReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sClusterRoleBindingService.UpdateClusterRoleBinding(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}

// DeleteClusterRoleBinding 删除集群角色绑定
// @Tags kubernetes
// @Summary 删除集群角色绑定
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "删除集群角色绑定"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /kubernetes/ClusterRoleBinding [delete]
func (k *K8sClusterRoleBindingApi) DeleteClusterRoleBinding(c *gin.Context) {
	req := clusterolebinding.DeleteClusterRoleBindingReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := k8sClusterRoleBindingService.DeleteClusterRoleBinding(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	} else {
		time.Sleep(1 * time.Second)
		response.OkWithMessage("删除成功", c)
	}
}

// CreateClusterRoleBinding 创建集群角色绑定
// @Tags kubernetes
// @Summary 创建集群角色绑定
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "创建集群角色绑定"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /kubernetes/ClusterRoleBinding [delete]
func (k *K8sClusterRoleBindingApi) CreateClusterRoleBinding(c *gin.Context) {
	req := clusterolebinding.CreateClusterRoleBindingReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if CronJob, err := k8sClusterRoleBindingService.CreateClusterRoleBinding(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(CronJob, c)
	}
}
