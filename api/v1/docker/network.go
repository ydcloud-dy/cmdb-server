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

type NetworkApi struct{}

// ListNetwork @Tags Docker
// @Summary 获取docker网络列表
// @Produce  application/json
// @Success 200
// @Router /docker/network/list [get]
func (p *NetworkApi) ListNetwork(c *gin.Context) {
	var plug model.SearchNetwork
	_ = c.ShouldBindQuery(&plug)

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}

	if res, err := service.ServiceGroupApp.ListNetwork(host, plug); err != nil {
		global.DYCLOUD_LOG.Error("获取网络列表失败!", zap.Error(err))
		response.FailWithMessage("获取网络列表失败", c)
	} else {
		networkRes := model.SearchNetworkRes{
			Pagination: model.Pagination{
				Page:     plug.Page,
				PageSize: plug.PageSize,
				Total:    len(res),
			},
		}
		if plug.PageSize == 0 && plug.Page == 0 {
			networkRes.Items = res
		} else {
			networkRes.Items = docker.SlicePagination(plug.Page, plug.PageSize, res)
		}
		response.OkWithDetailed(networkRes, "获取网络列表成功", c)
	}
}

// RemoveNetwork @Tags Docker
// @Summary 删除网络
// @Produce  application/json
// @Success 200
// @Router /docker/network [delete]
func (p *NetworkApi) RemoveNetwork(c *gin.Context) {

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}

	var plug model.RemoveNetwork
	_ = c.ShouldBindJSON(&plug)
	if err := service.ServiceGroupApp.RemoveNetwork(host, plug.Ids); err != nil {
		global.DYCLOUD_LOG.Error("删除网络失败!", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("删除网络失败,失败原因:%v", err.Error()), c)
	} else {
		response.Ok(c)
	}
}

// CreateNetwork @Tags Docker
// @Summary  创建网络
// @Produce  application/json
// @Success 200
// @Router /docker/network [post]
func (p *NetworkApi) CreateNetwork(c *gin.Context) {

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}
	var plug model.Network
	_ = c.ShouldBindJSON(&plug)
	if err := service.ServiceGroupApp.CreateNetwork(host, plug); err != nil {
		global.DYCLOUD_LOG.Error("创建网络失败!", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("创建网络失败:%v", err.Error()), c)
	} else {
		response.Ok(c)
	}
}
