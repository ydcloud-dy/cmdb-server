package api

import (
	"DYCLOUD/global"
	global2 "DYCLOUD/global/docker"
	"DYCLOUD/model/common/response"
	model "DYCLOUD/model/docker"
	service "DYCLOUD/service/docker"
	docker "DYCLOUD/utils/docker/docker"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ImageApi struct{}

// ListImage @Tags Docker
// @Summary 获取docker镜像列表
// @Produce  application/json
// @Success 200
// @Router /docker/image/list [get]
func (p *ImageApi) ListImage(c *gin.Context) {

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}

	var plug model.SearchImage
	_ = c.ShouldBindQuery(&plug)
	if res, err := service.ServiceGroupApp.ListImage(host, plug); err != nil {
		global.DYCLOUD_LOG.Error("获取镜像列表失败!", zap.Error(err))
		response.FailWithMessage("获取镜像列表失败", c)
	} else {
		imageRes := model.SearchImageRes{
			Pagination: model.Pagination{
				Page:     plug.Page,
				PageSize: plug.PageSize,
				Total:    len(res),
			},
		}
		if plug.PageSize == 0 && plug.Page == 0 {
			imageRes.Items = res
		} else {
			imageRes.Items = docker.SlicePagination(plug.Page, plug.PageSize, res)
		}
		response.OkWithDetailed(imageRes, "获取镜像列表成功", c)
	}
}

// PullImage @Tags Docker
// @Summary 下载镜像
// @Produce  application/json
// @Success 200
// @Router /docker/container/pull [post]
func (p *ImageApi) PullImage(c *gin.Context) {

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}
	var plug model.PullImage
	_ = c.ShouldBindJSON(&plug)
	if err := service.ServiceGroupApp.PullImage(host, plug); err != nil {
		global.DYCLOUD_LOG.Error("下载镜像失败!", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("下载镜像失败:"+err.Error()), c)
	} else {
		response.OkWithMessage("开始后台执行,下载镜像", c)
	}
}

// RemoveImage @Tags Docker
// @Summary 删除镜像
// @Produce  application/json
// @Success 200
// @Router /docker/container [delete]
func (p *ImageApi) RemoveImage(c *gin.Context) {

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}

	var plug model.RemoveImage
	_ = c.ShouldBindJSON(&plug)

	if err := service.ServiceGroupApp.RemoveImage(host, plug.Ids); err != nil {
		global.DYCLOUD_LOG.Error("删除容器失败!", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("删除容器失败,失败原因:%v", err.Error()), c)
	} else {
		response.Ok(c)
	}
}
