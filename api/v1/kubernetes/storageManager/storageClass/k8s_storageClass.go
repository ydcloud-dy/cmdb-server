package storageClass

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/model/kubernetes/storageClass"
	"DYCLOUD/service"
	"DYCLOUD/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type K8sStorageClassApi struct{}

var k8sStorageClassService = service.ServiceGroupApp.StorageClassServiceGroup.K8sStorageClassService

func (k *K8sStorageClassApi) GetStorageClassList(c *gin.Context) {
	req := storageClass.GetStorageClassListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := k8sStorageClassService.GetStorageClassList(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	} else {

		response.OkWithDetailed(storageClass.StorageClassListResponse{
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
func (k *K8sStorageClassApi) DescribeStorageClassInfo(c *gin.Context) {
	req := storageClass.DescribeStorageClassReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sStorageClassService.DescribeStorageClass(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(storageClass.DescribeStorageClassResponse{Items: list}, "获取成功", c)
	}
}
func (k *K8sStorageClassApi) UpdateStorageClass(c *gin.Context) {
	req := storageClass.UpdateStorageClassReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sStorageClassService.UpdateStorageClass(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}
func (k *K8sStorageClassApi) DeleteStorageClass(c *gin.Context) {
	req := storageClass.DeleteStorageClassReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := k8sStorageClassService.DeleteStorageClass(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	} else {
		time.Sleep(1 * time.Second)
		response.OkWithMessage("删除成功", c)
	}
}
func (k *K8sStorageClassApi) CreateStorageClass(c *gin.Context) {
	req := storageClass.CreateStorageClassReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if CronJob, err := k8sStorageClassService.CreateStorageClass(req, utils.GetUserUuid(c)); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(CronJob, c)
	}
}
