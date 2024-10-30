package namespaces

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	namespaces2 "DYCLOUD/model/kubernetes/namespaces"
	"DYCLOUD/service"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type K8sNamespaceApi struct{}

var k8sNamespaceService = service.ServiceGroupApp.NamespaceServiceGroup.K8sNamespaceService

func (k *K8sNamespaceApi) GetNamespaceList(c *gin.Context) {
	req := namespaces2.GetNamespaceListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := k8sNamespaceService.GetNamespaceList(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	} else {

		response.OkWithDetailed(namespaces2.NamespaceListResponse{
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
func (k *K8sNamespaceApi) DescribeNamespaceInfo(c *gin.Context) {
	req := namespaces2.DescribeNamespaceReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sNamespaceService.DescribeNamespace(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(namespaces2.DescribeNamespaceResponse{Items: list}, "获取成功", c)
	}
}
func (k *K8sNamespaceApi) UpdateNamespace(c *gin.Context) {
	req := namespaces2.UpdateNamespaceReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sNamespaceService.UpdateNamespace(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}
func (k *K8sNamespaceApi) DeleteNamespace(c *gin.Context) {
	req := namespaces2.DeleteNamespaceReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := k8sNamespaceService.DeleteNamespace(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	} else {
		time.Sleep(1 * time.Second)
		response.OkWithMessage("删除成功", c)
	}
}
func (k *K8sNamespaceApi) CreateNamespace(c *gin.Context) {
	req := namespaces2.CreateNamespaceReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if CronJob, err := k8sNamespaceService.CreateNamespace(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(CronJob, c)
	}
}
