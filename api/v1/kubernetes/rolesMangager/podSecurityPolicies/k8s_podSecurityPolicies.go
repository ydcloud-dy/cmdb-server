package podSecurityPolicies

import (
	"DYCLOUD/global"
	"DYCLOUD/model/PodSecurityPolicies"
	"DYCLOUD/model/common/request"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type K8sPodSecurityPoliciesApi struct{}

var k8sPodSecurityPoliciesService = service.ServiceGroupApp.PodSecurityPoliciesServiceGroup.K8sPodSecurityPoliciesService

func (k *K8sPodSecurityPoliciesApi) GetPodSecurityPoliciesList(c *gin.Context) {
	req := PodSecurityPolicies.GetPodSecurityPoliciesListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := k8sPodSecurityPoliciesService.GetPodSecurityPoliciesList(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	} else {

		response.OkWithDetailed(PodSecurityPolicies.PodSecurityPoliciesListResponse{
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
func (k *K8sPodSecurityPoliciesApi) DescribePodSecurityPoliciesInfo(c *gin.Context) {
	req := PodSecurityPolicies.DescribePodSecurityPoliciesReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sPodSecurityPoliciesService.DescribePodSecurityPolicies(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(PodSecurityPolicies.DescribePodSecurityPoliciesResponse{Items: list}, "获取成功", c)
	}
}
func (k *K8sPodSecurityPoliciesApi) UpdatePodSecurityPolicies(c *gin.Context) {
	req := PodSecurityPolicies.UpdatePodSecurityPoliciesReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sPodSecurityPoliciesService.UpdatePodSecurityPolicies(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}
func (k *K8sPodSecurityPoliciesApi) DeletePodSecurityPolicies(c *gin.Context) {
	req := PodSecurityPolicies.DeletePodSecurityPoliciesReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := k8sPodSecurityPoliciesService.DeletePodSecurityPolicies(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	} else {
		time.Sleep(1 * time.Second)
		response.OkWithMessage("删除成功", c)
	}
}
func (k *K8sPodSecurityPoliciesApi) CreatePodSecurityPolicies(c *gin.Context) {
	req := PodSecurityPolicies.CreatePodSecurityPoliciesReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if CronJob, err := k8sPodSecurityPoliciesService.CreatePodSecurityPolicies(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(CronJob, c)
	}
}
