package api

import (
	"DYCLOUD/global"
	global2 "DYCLOUD/global/docker"
	"DYCLOUD/model/common/response"
	model "DYCLOUD/model/docker"
	service "DYCLOUD/service/docker"
	"DYCLOUD/utils/docker/docker"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VolumeApi struct{}

// ListVolume @Tags Docker
// @Summary 获取存储卷列表
// @Produce  application/json
// @Success 200
// @Router /docker/volume/list [get]
func (p *VolumeApi) ListVolume(c *gin.Context) {

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}

	var plug model.SearchVolume
	_ = c.ShouldBindQuery(&plug)
	if res, err := service.ServiceGroupApp.ListVolume(host, plug); err != nil {
		global.DYCLOUD_LOG.Error("获取存储卷列表失败!", zap.Error(err))
		response.FailWithMessage("获取存储卷列表失败", c)
	} else {
		volumeRes := model.SearchVolumeRes{
			Pagination: model.Pagination{
				Page:     plug.Page,
				PageSize: plug.PageSize,
				Total:    len(res),
			},
		}
		if plug.PageSize == 0 && plug.Page == 0 {
			volumeRes.Items = res
		} else {
			volumeRes.Items = docker.SlicePagination(plug.Page, plug.PageSize, res)
		}
		response.OkWithDetailed(volumeRes, "获取存储卷列表成功", c)
	}
}

// RemoveVolume @Tags Docker
// @Summary 删除存储卷
// @Produce  application/json
// @Success 200
// @Router /docker/volume [delete]
func (p *VolumeApi) RemoveVolume(c *gin.Context) {

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}

	var plug model.RemoveNetwork
	_ = c.ShouldBindJSON(&plug)
	if err := service.ServiceGroupApp.RemoveVolume(host, plug.Ids); err != nil {
		global.DYCLOUD_LOG.Error("删除存储卷失败!", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("删除存储卷失败,失败原因:%v", err.Error()), c)
	} else {
		response.Ok(c)
	}
}

// CreateVolume @Tags Docker
// @Summary  创建存储卷
// @Produce  application/json
// @Success 200
// @Router /docker/volume [post]
func (p *VolumeApi) CreateVolume(c *gin.Context) {

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}

	var plug model.Volume
	_ = c.ShouldBindJSON(&plug)
	if err := service.ServiceGroupApp.CreateVolume(host, plug); err != nil {
		global.DYCLOUD_LOG.Error("创建存储卷失败!", zap.Error(err))
		response.FailWithMessage("创建存储卷失败", c)
	} else {
		response.Ok(c)
	}
}
