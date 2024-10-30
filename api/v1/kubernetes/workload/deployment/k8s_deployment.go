package deployment

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/model/kubernetes/deployment"
	"DYCLOUD/service"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type K8sDeploymentApi struct{}

var k8sDeploymentService = service.ServiceGroupApp.DeploymentServiceGroup.K8sDeploymentService

func (k *K8sDeploymentApi) GetDeploymentList(c *gin.Context) {
	var req = deployment.GetDeploymentListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := k8sDeploymentService.GetDeploymentList(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	} else {
		response.OkWithDetailed(deployment.DeploymentListResponse{
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

func (k *K8sDeploymentApi) DescribeDeploymentInfo(c *gin.Context) {
	req := deployment.DescribeDeploymentInfoReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deploy, err := k8sDeploymentService.DescribeDeploymentInfo(req, utils.GetUserUuid(c))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(deployment.DescribeDeployResponse{Items: deploy}, "获取成功", c)

}
func (k *K8sDeploymentApi) UpdateDeploymentInfo(c *gin.Context) {
	var req deployment.UpdateDeploymentInfoReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := k8sDeploymentService.UpdateDeploymentInfo(req, utils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		response.OkWithDetailed(data, "更新成功", c)
	}
}
func (k *K8sDeploymentApi) CreateDeployment(c *gin.Context) {
	var req deployment.CreateDeploymentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := k8sDeploymentService.CreateDeployment(req, utils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		response.OkWithDetailed(data, "更新成功", c)
	}
}
func (k *K8sDeploymentApi) DeleteDeployment(c *gin.Context) {
	var req deployment.DeleteDeploymentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = k8sDeploymentService.DeleteDeployment(req, utils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		response.OkWithMessage("删除成功", c)
	}
}
func (k *K8sDeploymentApi) RollBackDeployment(c *gin.Context) {
	var req deployment.RollBackDeployment
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := k8sDeploymentService.RollBackDeployment(req, utils.GetUserUuid(c))
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		response.OkWithDetailed(data, "更新成功", c)
	}
}
