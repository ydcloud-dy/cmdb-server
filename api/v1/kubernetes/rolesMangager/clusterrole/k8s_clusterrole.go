package clusterrole

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/model/kubernetes/clusterrole"
	"DYCLOUD/service"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type K8sClusterRoleApi struct{}

var k8sClusterRoleService = service.ServiceGroupApp.ClusterRoleServiceGroup.K8sClusterRoleService

// GetClusterRoleList 获取ClusterRole列表
// @Tags kubernetes
// @Summary 获取ClusterRole列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "获取ClusterRole列表"
// @Success 200 {object} response.Response{msg=string} "修改成功"
// @Router /kubernetes/ClusterRole [get]
func (k *K8sClusterRoleApi) GetClusterRoleList(c *gin.Context) {
	req := clusterrole.GetClusterRoleListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := k8sClusterRoleService.GetClusterRoleList(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	} else {

		response.OkWithDetailed(clusterrole.ClusterRoleListResponse{
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

// DescribeClusterRoleInfo 获取ClusterRole详情
// @Tags kubernetes
// @Summary 获取ClusterRole详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "获取ClusterRole详情"
// @Success 200 {object} response.Response{msg=string} "修改成功"
// @Router /kubernetes/ClusterRoleDetails [get]
func (k *K8sClusterRoleApi) DescribeClusterRoleInfo(c *gin.Context) {
	req := clusterrole.DescribeClusterRoleReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sClusterRoleService.DescribeClusterRole(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(clusterrole.DescribeClusterRoleResponse{Items: list}, "获取成功", c)
	}
}

// UpdateClusterRole 更新ClusterRole信息
// @Tags kubernetes
// @Summary 更新ClusterRole信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "更新ClusterRole信息"
// @Success 200 {object} response.Response{msg=string} "修改成功"
// @Router /kubernetes/ClusterRole [put]
func (k *K8sClusterRoleApi) UpdateClusterRole(c *gin.Context) {
	req := clusterrole.UpdateClusterRoleReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sClusterRoleService.UpdateClusterRole(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}

// DeleteClusterRole 删除ClusterRole
// @Tags kubernetes
// @Summary 删除ClusterRole
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "删除ClusterRole"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /kubernetes/ClusterRole [put]
func (k *K8sClusterRoleApi) DeleteClusterRole(c *gin.Context) {
	req := clusterrole.DeleteClusterRoleReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := k8sClusterRoleService.DeleteClusterRole(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	} else {
		time.Sleep(1 * time.Second)
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteClusterRole 创建集群角色
// @Tags kubernetes
// @Summary 创建集群角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "创建集群角色"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /kubernetes/ClusterRole [post]
func (k *K8sClusterRoleApi) CreateClusterRole(c *gin.Context) {
	req := clusterrole.CreateClusterRoleReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if CronJob, err := k8sClusterRoleService.CreateClusterRole(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(CronJob, c)
	}
}
