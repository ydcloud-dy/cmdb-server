package api

import (
	"DYCLOUD/global"
	global2 "DYCLOUD/global/docker"
	"DYCLOUD/model/common/response"
	model "DYCLOUD/model/docker"
	service "DYCLOUD/service/docker"
	docker "DYCLOUD/utils/docker/docker"
	"DYCLOUD/utils/docker/str_util"
	"bufio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type ContainerApi struct{}

// ListContainer @Tags Docker
// @Summary 获取docker容器列表
// @Produce  application/json
// @Success 200
// @Router /docker/container/list [get]
func (p *ContainerApi) ListContainer(c *gin.Context) {

	var plug model.SearchContainer
	_ = c.ShouldBindQuery(&plug)

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}
	if res, err := service.ServiceGroupApp.ListContainer(host, plug); err != nil {
		global.DYCLOUD_LOG.Error("获取容器列表失败:", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取容器列表失败:%v", err.Error()), c)
	} else {
		containerRes := model.SearchContainerRes{
			Pagination: model.Pagination{
				Page:     plug.Page,
				PageSize: plug.PageSize,
				Total:    len(res),
			},
		}
		if plug.PageSize == 0 && plug.Page == 0 {
			containerRes.Items = res
		} else {
			containerRes.Items = docker.SlicePagination(plug.Page, plug.PageSize, res)
		}
		response.OkWithDetailed(containerRes, "获取容器列表成功", c)
	}
}

// AddContainer @Tags Docker
// @Summary 创建容器
// @Produce  application/json
// @Success 200
// @Router /container [post]
func (p *ContainerApi) AddContainer(c *gin.Context) {

	var plug model.AddContainer
	_ = c.ShouldBindJSON(&plug)

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}
	if err := service.ServiceGroupApp.AddContainer(host, plug); err != nil {
		global.DYCLOUD_LOG.Error("创建容器失败!", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("创建容器失败:%v", err.Error()), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// EditContainer @Tags Docker
// @Summary 修改容器
// @Produce  application/json
// @Success 200
// @Router /docker/container [put]
func (p *ContainerApi) EditContainer(c *gin.Context) {

	var plug model.UpdateContainer
	_ = c.ShouldBindJSON(&plug)

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}
	if err := service.ServiceGroupApp.UpdateContainer(host, plug); err != nil {
		global.DYCLOUD_LOG.Error("修改容器失败:", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("修改容器失败:%v", err.Error()), c)
	} else {
		response.OkWithMessage("修改成功", c)

	}
}

// RemoveContainer @Tags Docker
// @Summary 删除容器
// @Produce  application/json
// @Success 200
// @Router /docker/container [delete]
func (p *ContainerApi) RemoveContainer(c *gin.Context) {

	var plug model.RemoveContainer
	_ = c.ShouldBindJSON(&plug)
	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}
	if err := service.ServiceGroupApp.RemoveContainer(host, plug); err != nil {
		global.DYCLOUD_LOG.Error("删除容器失败:", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("删除容器失败:%v", err.Error()), c)
	} else {
		response.OkWithMessage("删除容器成功", c)
	}
}

// RestartContainer @Tags Docker
// @Summary 重启容器
// @Produce  application/json
// @Success 200
// @Router /docker/container/restart [post]
func (p *ContainerApi) RestartContainer(c *gin.Context) {
	var plug model.RestartContainer
	_ = c.ShouldBindJSON(&plug)
	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}
	if err := service.ServiceGroupApp.RestartContainer(host, plug); err != nil {
		global.DYCLOUD_LOG.Error("重启容器失败!", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("重启容器失败:%v", err.Error()), c)
	} else {
		response.OkWithMessage("重启容器成功", c)
	}
}

// StopContainer @Tags Docker
// @Summary 停止容器
// @Produce  application/json
// @Success 200
// @Router /docker/container/stop [post]
func (p *ContainerApi) StopContainer(c *gin.Context) {
	var plug model.StopContainer
	_ = c.ShouldBindJSON(&plug)
	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}
	if err := service.ServiceGroupApp.StopContainer(host, plug); err != nil {
		global.DYCLOUD_LOG.Error("停止容器失败!", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("停止容器失败:%v", err.Error()), c)
	} else {
		response.OkWithMessage("停止容器成功", c)
	}
}

// StartContainer @Tags Docker
// @Summary 启动容器
// @Produce  application/json
// @Success 200
// @Router /docker/container/start [post]
func (p *ContainerApi) StartContainer(c *gin.Context) {
	var plug model.StartContainer
	_ = c.ShouldBindJSON(&plug)

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}
	if err := service.ServiceGroupApp.StartContainer(host, plug); err != nil {
		global.DYCLOUD_LOG.Error("启动容器失败!", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("启动容器失败:%v", err.Error()), c)
	} else {
		response.OkWithMessage("启动容器成功", c)
	}
}

// StatsContainer @Tags Docker
// @Summary 容器状态监控
// @Produce  application/json
// @Success 200
// @Router /docker/container/stats [post]
func (p *ContainerApi) StatsContainer(c *gin.Context) {
	id := c.Query("id")
	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}
	if stats, err := service.ServiceGroupApp.GetContainerStats(host, model.StatsContainer{Id: id}); err != nil {
		global.DYCLOUD_LOG.Error("容器监控指标失败!", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("容器监控指标失败:%v", err.Error()), c)
	} else {
		response.OkWithData(stats, c)
	}
}

// LogsContainer @Tags Docker
// @Summary 查看容器日志
// @Produce  application/json
// @Success 200
// @Router /docker/container/log [get]
func (p *ContainerApi) LogsContainer(c *gin.Context) {
	plug := &model.GetContainerLog{}
	err := c.ShouldBindQuery(plug)
	if err != nil {
		return
	}

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}
	cLogRes, err := service.ServiceGroupApp.GetContainerLog(host, plug)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("查看容器日志失败:%v", err.Error()), c)
		return
	}

	scanner := bufio.NewScanner(cLogRes.Reader)
	buf := new(bytes.Buffer)
	var val, text, line string
	for scanner.Scan() {
		b := scanner.Bytes()
		if len(b) > 8 {
			line = string(b[8:])
		} else {
			line = string(b)
		}

		if plug.Text != "" {
			val = line
			text = plug.Text
			if plug.Regexp {
				if plug.UpperLower {
					text = "(?i)" + plug.Text
				}
				if res := str_util.FindRegex(val, text); res == "" {
					continue
				}
			} else {
				if s := str_util.ReplaceIgnoreCaseKeepCaseWithWrapper(plug.UpperLower, val, text, text, "<span style='color:yellow'>", "</span>"); s == "" {
					continue
				} else {
					line = s
				}
			}
		}
		buf.WriteString(line + "\n")
	}
	io.Copy(c.Writer, buf)
}

// ExecContainer @Tags Docker
// @Summary 开启容器终端
// @Produce  application/json
// @Success 200
// @Router /docker/container/exec [post]
func (p *ContainerApi) ExecContainer(c *gin.Context) {

	var plug model.ExecContainer
	_ = c.ShouldBindJSON(&plug)

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}
	if idResponse, err := service.ServiceGroupApp.ExecContainer(host, plug); err != nil {
		global.DYCLOUD_LOG.Error("开启容器终端:", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("开启容器终端:%v", err), c)
	} else {
		response.OkWithData(idResponse, c)
	}
}

// ExecContainerResize @Tags Docker
// @Summary 修改容器终端大小
// @Produce  application/json
// @Success 200
// @Router /docker/container/exec/resize [post]
func (p *ContainerApi) ExecContainerResize(c *gin.Context) {

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}
	var plug model.ExecContainerResize
	_ = c.ShouldBindJSON(&plug)
	_ = service.ServiceGroupApp.ExecContainerResize(host, plug)
	response.OkWithMessage("操作成功", c)
}

// UpGrader 创建UpGrader对象
var UpGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ExecTermContainer @Tags Docker
// @Summary 终端执行命令websocket
// @Produce  application/json
// @Success 200
// @Router /container/exec/term/:id[websocket]
func (p *ContainerApi) ExecTermContainer(c *gin.Context) {

	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}
	conn, err := UpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	id := c.Param("id")
	targetConn, buf, err := service.ServiceGroupApp.ExecAttachContainer(host, id)
	if err != nil {
		return
	}
	go func() {
		b := [1024]byte{}
		for {
			n, err := buf.Read(b[:])
			if err != nil {
				break
			}
			if n <= 0 {
				continue
			}
			_ = conn.WriteMessage(websocket.BinaryMessage, b[:n])
		}
	}()
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			_, _ = targetConn.Write(message)
		}
	}()
}

// InspectContainer @Tags Docker
// @Summary 查看容器详细信息
// @Produce  application/json
// @Success 200
// @Router /docker/container/inspect [post]
func (p *ContainerApi) InspectContainer(c *gin.Context) {
	var plug model.InspectContainer
	_ = c.ShouldBindJSON(&plug)
	host, err := global2.Context.GetHost(c)
	if err != nil {
		response.FailWithMessage("主机名称格式错误!", c)
		return
	}
	if res, err := service.ServiceGroupApp.InspectContainer(host, plug); err != nil {
		global.DYCLOUD_LOG.Error("获取容器详细信息失败!", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取容器详细信息失败:%v", err.Error()), c)
	} else {
		response.OkWithData(res, c)
	}
}
