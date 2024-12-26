package ws

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/response"
	kubernetesReq "DYCLOUD/model/kubernetes/ws" // 确保导入路径正确
	"DYCLOUD/model/system/request"
	"DYCLOUD/utils"
	sutils "DYCLOUD/utils"
	"DYCLOUD/utils/kubernetes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type WsApi struct{}

// 终端token验证
func GetJWTAuth(token string) (user *request.CustomClaims, err error) {
	j := sutils.NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.DYCLOUD_LOG.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构")
	}

	return claims, err
}

// @Tags WsApi
// @Summary  终端
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query kubernetesReq.TerminalRequest true "终端请求参数"
// @Success 200 {object} response.Response{msg=string} "终端请求成功，返回包括api详情"
// @Router /kubernetes/pods/terminal [get]
func (w *WsApi) Terminal(c *gin.Context) {
	var terminal kubernetesReq.TerminalRequest
	_ = c.ShouldBindQuery(&terminal)
	user, err := GetJWTAuth(terminal.XToken)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}

	if err := sutils.Verify(terminal, utils.TerminalVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	client, err := kubernetes.NewKubeClient(terminal.ClusterId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("连接k8s client失败!，%v", err.Error()), c)
	}

	kubeshell, err := kubernetes.NewKubeShell(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	cmd := []string{
		"/bin/sh", "-c", fmt.Sprintf("clear;(bash || sh); export LINES=%d ; export COLUMNS=%d;", terminal.Rows, terminal.Cols),
	}

	// 创建一个通道用于通知退出
	exitChan := make(chan bool)

	// 超时退出websocket
	go func() {
		for {
			select {
			case <-exitChan:
				global.DYCLOUD_LOG.Error("终端记录已经写入退出循环.")
				return // 接收到退出通知时退出循环
			case <-time.After(time.Second * 10):
				if time.Now().Unix()-kubeshell.UpdatedAt.Unix() > 120 || time.Now().Unix()-kubeshell.WriteAt.Unix() > 120 {
					if _, err = kubeshell.Write([]byte("exit\n")); err != nil {
						global.DYCLOUD_LOG.Error("终端WebSocket 终端退出失败：" + err.Error())
					}

					if err := kubeshell.WriteLog(terminal, user); err != nil {
						global.DYCLOUD_LOG.Error("终端记录写入失败：" + err.Error())
						return
					}

					if err = kubeshell.Close(); err != nil {
						global.DYCLOUD_LOG.Error("终端WebSocket超时关闭失败：" + err.Error())
						return
					}

					global.DYCLOUD_LOG.Info("终端WebSocket超时关闭成功！")
					return
				}
			}
		}
	}()

	if err := client.Pod.Exec(cmd, kubeshell, terminal); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := kubeshell.WriteLog(terminal, user); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// 通知进程退出循环
	exitChan <- true
}

// @Tags WsApi
// @Summary  终端日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query kubernetesReq.TerminalRequest true "终端请求参数"
// @Success 200 {object} response.Response{msg=string} "终端请求成功，返回包括api详情"
// @Router /kubernetes/pods/logs [get]
func (w *WsApi) ContainerLog(c *gin.Context) {
	var terminal kubernetesReq.TerminalRequest
	_ = c.ShouldBindQuery(&terminal)
	if _, err := GetJWTAuth(terminal.XToken); err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	if err := sutils.Verify(terminal, utils.TerminalVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	kubeLogger, err := kubernetes.NewKubeLogger(c.Writer, c.Request, nil)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("连接终端websocket升级失败!，%v", err.Error()), c)
	}

	client, err := kubernetes.NewKubeClient(terminal.ClusterId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("连接k8s client失败!，%v", err.Error()), c)
	}

	client.Pod.ContainerLog(kubeLogger, terminal.Name, terminal.PodName, terminal.Namespace)
}
