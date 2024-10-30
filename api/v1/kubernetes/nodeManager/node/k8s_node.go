package node

import (
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/model/kubernetes"
	"DYCLOUD/model/kubernetes/nodes"
	"DYCLOUD/service"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
)

type K8sNodeApi struct{}

var k8sNodeService = service.ServiceGroupApp.NodeServiceGroup.K8sNodeService

// GetNodeList 获取node列表
// @Tags kubernetes
// @Summary 获取node列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "获取node列表"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /kubernetes/nodes [get]
func (k *K8sNodeApi) GetNodeList(c *gin.Context) {
	req := nodes.NodeListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(req, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	nodeList, total, err := k8sNodeService.GetNodeList(req, utils.GetUserUuid(c))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(nodes.NodeListResponse{
		Items: nodeList,
		Total: total,
		PageInfo: request.PageInfo{
			Page:     req.Page,
			PageSize: req.PageSize,
			Keyword:  req.Keyword,
		},
	}, "获取成功", c)
}

// GetNodeMetricsList 获取node监控指标
// @Tags kubernetes
// @Summary 获取node监控指标
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "获取node监控指标"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /kubernetes/nodes/metrics [get]
func (k *K8sNodeApi) GetNodeMetricsList(c *gin.Context) {
	req := nodes.NodeMetricsReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(req, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	nodeMetricsList, err := k8sNodeService.GetNodeMetricsList(req, utils.GetUserUuid(c))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//response.OkWithDetailed(nodeMetricsList, "获取成功", c)
	response.OkWithData(nodeMetricsList, c)
}

// GetNodeMetricsList 获取node详情
// @Tags kubernetes
// @Summary 获取node详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "获取node详情"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /kubernetes/nodeDetails [get]
func (k *K8sNodeApi) DescribeNodeInfo(c *gin.Context) {
	req := nodes.DescribeNodeReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(req, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	nodeInfo, err := k8sNodeService.DescribeNodeInfo(req, utils.GetUserUuid(c))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//response.OkWithDetailed(nodeMetricsList, "获取成功", c)
	response.OkWithData(nodes.DescribeNodeInfoResponse{Items: nodeInfo}, c)
}

// UpdateNodeInfo 修改node信息
// @Tags kubernetes
// @Summary 修改node信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "修改node信息"
// @Success 200 {object} response.Response{msg=string} "修改成功"
// @Router /kubernetes/nodes [get]
func (k *K8sNodeApi) UpdateNodeInfo(c *gin.Context) {
	req := nodes.UpdateNodeReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(req, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	nodeInfo, err := k8sNodeService.UpdateNodeInfo(req, utils.GetUserUuid(c))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//response.OkWithDetailed(nodeMetricsList, "获取成功", c)
	response.OkWithData(nodes.DescribeNodeInfoResponse{Items: nodeInfo}, c)
}

// EvictAllNodePod 驱逐node上所有的pod
// @Tags kubernetes
// @Summary 驱逐node上所有的pod
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "驱逐node上所有的pod"
// @Success 200 {object} response.Response{msg=string} "修改成功"
// @Router /kubernetes/nodes/EvictAllPod [post]
func (k *K8sNodeApi) EvictAllNodePod(c *gin.Context) {
	req := nodes.EvictAllNodePodReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(req, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = k8sNodeService.EvictAllNodePod(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//response.OkWithDetailed(nodeMetricsList, "获取成功", c)
	response.OkWithMessage("操作成功", c)
}
