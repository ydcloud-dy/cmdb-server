package cluster

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/model/kubernetes"
	"DYCLOUD/model/kubernetes/cluster"
	clusterReq "DYCLOUD/model/kubernetes/cluster/request"
	response2 "DYCLOUD/model/kubernetes/cluster/response"
	"DYCLOUD/service"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type K8sClusterApi struct{}

var k8sClusterService = service.ServiceGroupApp.ClusterServiceGroup.K8sClusterService

// CreateK8sCluster 创建k8sCluster表
// @Tags K8sCluster
// @Summary 创建k8sCluster表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "创建k8sCluster表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /k8sCluster/createK8sCluster [post]
func (k8sClusterApi *K8sClusterApi) CreateK8sCluster(c *gin.Context) {
	var k8sCluster = cluster.K8sCluster{}
	err := c.ShouldBindJSON(&k8sCluster)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	k8sCluster.CreatedBy = utils.GetUserID(c)

	if err := k8sClusterService.CreateK8sCluster(&k8sCluster); err != nil {
		global.DYCLOUD_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteK8sCluster 删除k8sCluster表
// @Tags K8sCluster
// @Summary 删除k8sCluster表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "删除k8sCluster表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /k8sCluster/deleteK8sCluster [delete]
func (k8sClusterApi *K8sClusterApi) DeleteK8sCluster(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := k8sClusterService.DeleteK8sCluster(idInfo.ID); err != nil {
		global.DYCLOUD_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteK8sClusterByIds 批量删除k8sCluster表
// @Tags K8sCluster
// @Summary 批量删除k8sCluster表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /k8sCluster/deleteK8sClusterByIds [delete]
func (k8sClusterApi *K8sClusterApi) DeleteK8sClusterByIds(c *gin.Context) {
	var ids request.IdsReq
	_ = c.ShouldBindJSON(&ids)
	if err := k8sClusterService.DeleteK8sClusterByIds(ids); err != nil {
		global.DYCLOUD_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateK8sCluster 更新k8sCluster表
// @Tags K8sCluster
// @Summary 更新k8sCluster表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "更新k8sCluster表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /k8sCluster/updateK8sCluster [put]
func (k8sClusterApi *K8sClusterApi) UpdateK8sCluster(c *gin.Context) {
	var k8sCluster cluster.K8sCluster
	err := c.ShouldBindJSON(&k8sCluster)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	k8sCluster.UpdatedBy = utils.GetUserID(c)

	if err := k8sClusterService.UpdateK8sCluster(k8sCluster); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// CreateCredential 创建集群凭据
// @Tags K8sCluster
// @Summary 创建集群凭据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "创建集群凭据"
// @Success 200 {object} response.Response{msg=string} "凭据创建成功"
// @Router /k8sCluster/CreateCredential [post]
func (K8sClusterApi *K8sClusterApi) CreateCredential(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := k8sClusterService.CreateCredential(idInfo.ID, utils.GetUserID(c)); err != nil {
		global.DYCLOUD_LOG.Error("凭据创建失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("凭据创建成功", c)
	}
}

// GetClusterUserById 根据id获取集群用户
// @Tags K8sCluster
// @Summary 根据id获取集群用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "根据id获取集群用户"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /k8sCluster/getUserById [get]
func (K8sClusterApi *K8sClusterApi) GetClusterUserById(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cluster, err := k8sClusterService.GetClusterUserById(idInfo.ID, utils.GetUserID(c))
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response2.ClusterUserResponse{User: cluster}, "获取成功", c)
	}
}

// GetClusterRoles 获取集群角色
// @Tags K8sCluster
// @Summary 获取集群角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "获取集群角色"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /k8sCluster/getClusterRoles [get]
func (K8sClusterApi *K8sClusterApi) GetClusterRoles(c *gin.Context) {
	var roleType clusterReq.ClusterRoleType
	_ = c.ShouldBindJSON(&roleType)
	if err := utils.Verify(roleType, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	roles, err := k8sClusterService.GetClusterRoles(roleType)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response2.RolesResponse{Roles: roles}, "获取成功", c)
	}
}

// GetClusterApiGroups 获取集群api组
// @Tags K8sCluster
// @Summary 获取集群api组
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "获取集群api组"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /k8sCluster/getClusterApiGroups [get]
func (K8sClusterApi *K8sClusterApi) GetClusterApiGroups(c *gin.Context) {
	var apiGroups clusterReq.ClusterApiGroups
	_ = c.ShouldBindJSON(&apiGroups)
	if err := utils.Verify(apiGroups, kubernetes.ApiGroupsVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	groups, err := k8sClusterService.GetClusterApiGroups(apiGroups)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response2.ApiGroupResponse{Groups: groups}, "获取成功", c)
	}
}

// CreateClusterRole 创建集群角色
// @Tags K8sCluster
// @Summary 创建集群角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "创建集群角色"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /k8sCluster/createClusterRole [post]
func (K8sClusterApi *K8sClusterApi) CreateClusterRole(c *gin.Context) {
	var role cluster.RoleData
	_ = c.ShouldBindJSON(&role)
	if err := utils.Verify(role, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := k8sClusterService.CreateClusterRole(role)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("角色创建成功!", c)
	}
}

// UpdateClusterRole 更新集群角色
// @Tags K8sCluster
// @Summary 更新集群角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "更新集群角色"
// @Success 200 {object} response.Response{msg=string} "角色更新成功"
// @Router /k8sCluster/updateClusterRole [put]
func (K8sClusterApi *K8sClusterApi) UpdateClusterRole(c *gin.Context) {
	var role cluster.RoleData
	_ = c.ShouldBindJSON(&role)
	if err := utils.Verify(role, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := k8sClusterService.UpdateClusterRole(role)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("角色更新成功!", c)
	}
}

// DeleteClusterRole 删除集群角色
// @Tags K8sCluster
// @Summary 删除集群角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "删除集群角色"
// @Success 200 {object} response.Response{msg=string} "角色删除成功"
// @Router /k8sCluster/deleteClusterRole [put]
func (K8sClusterApi *K8sClusterApi) DeleteClusterRole(c *gin.Context) {
	var role cluster.RoleData
	_ = c.ShouldBindJSON(&role)
	if err := utils.Verify(role, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := k8sClusterService.DeleteClusterRole(role)
	if err != nil {
		global.DYCLOUD_LOG.Error("角色删除失败!", zap.Error(err))
		response.FailWithMessage("角色删除失败!", c)
	} else {
		response.OkWithMessage("角色删除成功!", c)
	}
}

// CreateClusterUser 创建集群用户
// @Tags K8sCluster
// @Summary 创建集群用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "创建集群用户"
// @Success 200 {object} response.Response{msg=string} "用户授权成功"
// @Router /k8sCluster/createClusterUser [post]
func (K8sClusterApi *K8sClusterApi) CreateClusterUser(c *gin.Context) {
	var role clusterReq.CreateClusterRole
	_ = c.ShouldBindJSON(&role)
	if err := utils.Verify(role, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := k8sClusterService.CreateUser(role)
	if err != nil {
		global.DYCLOUD_LOG.Error("用户授权失败!", zap.Error(err))
		response.FailWithMessage("用户授权失败!", c)
	} else {
		response.OkWithMessage("用户授权成功!", c)
	}
}

// UpdateClusterUser 更新集群用户授权信息
// @Tags K8sCluster
// @Summary 更新集群用户授权信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "更新集群用户授权信息"
// @Success 200 {object} response.Response{msg=string} "用户授权更新成功"
// @Router /k8sCluster/updateClusterUser [put]
func (K8sClusterApi *K8sClusterApi) UpdateClusterUser(c *gin.Context) {
	var role clusterReq.CreateClusterRole
	_ = c.ShouldBindJSON(&role)
	if err := utils.Verify(role, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := k8sClusterService.UpdateClusterUser(role)
	if err != nil {
		global.DYCLOUD_LOG.Error("用户授权更新失败!", zap.Error(err))
		response.FailWithMessage("用户授权更新失败!", c)
	} else {
		response.OkWithMessage("用户授权更新成功!", c)
	}
}

// DeleteClusterUser 删除集群用户
// @Tags K8sCluster
// @Summary 删除集群用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "删除集群用户"
// @Success 200 {object} response.Response{msg=string} "用户删除成功"
// @Router /k8sCluster/deleteClusterUser [delete]
func (K8sClusterApi *K8sClusterApi) DeleteClusterUser(c *gin.Context) {
	var role clusterReq.DeleteClusterRole
	_ = c.ShouldBindJSON(&role)
	if err := utils.Verify(role, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := k8sClusterService.DeleteClusterUser(role)
	if err != nil {
		global.DYCLOUD_LOG.Error("用户删除失败!", zap.Error(err))
		response.FailWithMessage("用户删除失败!", c)
	} else {
		response.OkWithMessage("用户删除成功!", c)
	}
}

// GetClusterUserNamespace 获取集群用户命名空间
// @Tags K8sCluster
// @Summary 获取集群用户命名空间
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "获取集群用户命名空间"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /k8sCluster/getClusterUserNamespace [get]
func (K8sClusterApi *K8sClusterApi) GetClusterUserNamespace(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	namespaces, err := k8sClusterService.GetClusterUserNamespace(idInfo.ID, utils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error("获取命名空间失败!", zap.Error(err))
		response.FailWithMessage("获取命名空间失败", c)
	} else {
		response.OkWithDetailed(response2.ClusterUserNamespace{Namespaces: namespaces}, "获取成功", c)
	}
}

// GetClusterListNamespace 获取集群所有namespace
// @Tags K8sCluster
// @Summary 获取集群所有namespace
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cluster.K8sCluster true "获取集群所有namespace"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /k8sCluster/getClusterListNamespace [get]
func (k8sClusterApi *K8sClusterApi) GetClusterListNamespace(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	namespaces, err := k8sClusterService.GetClusterListNamespace(idInfo.ID)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取命名空间失败!", zap.Error(err))
		response.FailWithMessage("获取命名空间失败", c)
	} else {
		response.OkWithDetailed(response2.ClusterListNamespace{Namespaces: namespaces}, "获取成功", c)
	}
}

// FindK8sCluster 用id查询k8sCluster表
// @Tags K8sCluster
// @Summary 用id查询k8sCluster表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query cluster.K8sCluster true "用id查询k8sCluster表"
// @Success 200 {object} response.Response{data=object{rek8sCluster=cluster.K8sCluster},msg=string} "查询成功"
// @Router /k8sCluster/findK8sCluster [get]
func (k8sClusterApi *K8sClusterApi) FindK8sCluster(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cluster, err := k8sClusterService.GetK8sCluster(idInfo.ID)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response2.ClusterResponse{Cluster: cluster}, "获取成功", c)
	}
}
func (k8sClusterApi *K8sClusterApi) FindK8sClusterByName(c *gin.Context) {
	//var idInfo request.GetById
	//_ = c.ShouldBindJSON(&idInfo)

	//if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	name := c.Query("name")
	cluster, err := k8sClusterService.GetK8sClusterByName(name)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response2.ClusterResponse{Cluster: cluster}, "获取成功", c)
	}
}

// GetK8sClusterList 分页获取k8sCluster表列表
// @Tags K8sCluster
// @Summary 分页获取k8sCluster表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clusterReq.K8sClusterSearch true "分页获取k8sCluster表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /k8sCluster/getK8sClusterList [get]
func (k8sClusterApi *K8sClusterApi) GetK8sClusterList(c *gin.Context) {
	var pageInfo clusterReq.K8sClusterSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := k8sClusterService.GetK8sClusterInfoList(pageInfo); err != nil {
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

// GetK8sClusterPublic 不需要鉴权的k8sCluster表接口
// @Tags K8sCluster
// @Summary 不需要鉴权的k8sCluster表接口
// @accept application/json
// @Produce application/json
// @Param data query clusterReq.K8sClusterSearch true "分页获取k8sCluster表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /k8sCluster/getK8sClusterPublic [get]
func (k8sClusterApi *K8sClusterApi) GetK8sClusterPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的k8sCluster表接口信息",
	}, "获取成功", c)
}
