package statefulSet

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/model/kubernetes/statefulSet"
	"DYCLOUD/service"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type K8sStatefulSetApi struct{}

var k8sStatefulSetService = service.ServiceGroupApp.StatefulSetServiceGroup.K8sStatefulSetService

func (k *K8sStatefulSetApi) CreateStatefulset(c *gin.Context) {
	req := statefulSet.CreateStatefulSetReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if daemonset, err := k8sStatefulSetService.CreateStatefulSet(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(daemonset, c)
	}
}
func (k *K8sStatefulSetApi) UpdateStatefulSetInfo(c *gin.Context) {
	req := statefulSet.UpdateStatefulSetReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sStatefulSetService.UpdateStatefulSet(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}
func (k *K8sStatefulSetApi) GetStatefulSetList(c *gin.Context) {
	router := statefulSet.GetStatefulSetListReq{}
	req := router
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := k8sStatefulSetService.GetStatefulSetList(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(statefulSet.StatefulSetListResponse{
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

func (k *K8sStatefulSetApi) DescribeStatefulSetInfo(c *gin.Context) {
	req := statefulSet.DescribeStatefulSetReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sStatefulSetService.DescribeStatefulSet(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(statefulSet.DescribeStatefulSetResponse{Items: list}, "获取成功", c)
	}
}

func (k *K8sStatefulSetApi) DeleteStatefulSet(c *gin.Context) {
	req := statefulSet.DeleteStatefulSetReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := k8sStatefulSetService.DeleteStatefulSet(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	} else {
		response.OkWithMessage("删除成功", c)
	}
}
