package pod

import (
	"DYCLOUD/api/v1/ws"
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/model/kubernetes"
	"DYCLOUD/model/kubernetes/pods"
	"DYCLOUD/service"
	"DYCLOUD/utils"
	"archive/tar"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

type K8sPodApi struct{}

var k8sPodService = service.ServiceGroupApp.PodServiceGroup.K8sPodService

// GetPodList 获取 Pod 列表
// @Tags kubernetes
// @Summary 获取 Pod 列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query pods.PodListReq true "获取 Pod 列表"
// @Success 200 {object} response.Response{data=pods.PodListResponse,msg=string} "获取成功"
// @Router /kubernetes/pods [get]
func (k *K8sPodApi) GetPodList(c *gin.Context) {
	req := pods.PodListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(req, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	podList, total, err := k8sPodService.GetPodList(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(pods.PodListResponse{
		Items: podList,
		Total: total,
		PageInfo: request.PageInfo{
			Page:     req.Page,
			PageSize: req.PageSize,
			Keyword:  req.Keyword,
		},
	}, "获取成功", c)
}

// MetricsPodsList 获取 Pod Metrics 列表
// @Tags kubernetes
// @Summary 获取 Pod Metrics 列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query pods.PodMetricsReq true "获取 Pod Metrics 列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /kubernetes/pods/metrics [get]
func (k *K8sPodApi) MetricsPodsList(c *gin.Context) {
	req := pods.PodMetricsReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(req, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	nodeMetricsList, err := k8sPodService.GetPodMetricsList(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(nodeMetricsList, "获取成功", c)
}

// DescribePodInfo 获取 Pod 详情
// @Tags kubernetes
// @Summary 获取 Pod 详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query pods.DescribePodInfo true "获取 Pod 详情"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /kubernetes/podDetails [get]
func (k *K8sPodApi) DescribePodInfo(c *gin.Context) {
	req := pods.DescribePodInfo{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(req, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	podInfo, err := k8sPodService.DescribePodInfo(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(pods.DescribePodInfoResponse{Items: podInfo}, "获取成功", c)
}

// PodEvents 获取 Pod 事件列表
// @Tags kubernetes
// @Summary 获取 Pod 事件列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query pods.PodEventsReq true "获取 Pod 事件列表"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /kubernetes/events [get]
func (k *K8sPodApi) PodEvents(c *gin.Context) {
	req := pods.PodEventsReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(req, kubernetes.RoleTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	eventList, total, err := k8sPodService.PodEvents(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(pods.EventInfoResponse{Items: eventList, Total: total}, "获取成功", c)
}

// DownloadFile 下载 Pod 文件
// @Tags kubernetes
// @Summary 下载 Pod 文件
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query pods.PodsFilesRequest true "下载 Pod 文件"
// @Success 200 {object} response.Response{msg=string} "下载成功"
// @Router /kubernetes/pods/downloadFile [get]
func (p *K8sPodApi) DownloadFile(c *gin.Context) {
	var pods pods.PodsFilesRequest
	_ = c.ShouldBindQuery(&pods)

	if _, err := ws.GetJWTAuth(pods.XToken); err != nil {
		response.FailWithMessage(err.Error(), c)
	}

	if err := utils.Verify(pods, utils.PodVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	file, err := k8sPodService.DownloadFile(pods)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
	} else {
		c.Header("Content-Disposition", "attachment; filename="+path.Base(file))
		c.Header("Content-Type", "application/octet-stream")
		c.File(file)
	}
}

// CreatePod 创建 Pod
// @Tags kubernetes
// @Summary 创建 Pod
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body pods.CreatePodReq true "创建 Pod"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /kubernetes/pods [post]
func (k *K8sPodApi) CreatePod(c *gin.Context) {
	req := pods.CreatePodReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if pod, err := k8sPodService.CreatePod(req); err != nil {
		global.DYCLOUD_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(pod, c)
	}
}

// DeletePod 删除 Pod
// @Tags kubernetes
// @Summary 删除 Pod
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body pods.DeletePodReq true "删除 Pod"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /kubernetes/pods [delete]
func (k *K8sPodApi) DeletePod(c *gin.Context) {
	req := pods.DeletePodReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err = k8sPodService.DeletePod(req); err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败："+err.Error(), c)
		return
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// UpdatePod 更新 Pod
// @Tags kubernetes
// @Summary 更新 Pod
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body pods.UpdatePodReq true "更新 Pod"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /kubernetes/pods [put]
func (k *K8sPodApi) UpdatePod(c *gin.Context) {
	req := pods.UpdatePodReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if pod, err := k8sPodService.UpdatePod(req); err != nil {
		global.DYCLOUD_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(pod, c)
	}
}

// ListPodFiles 获取 Pod 文件列表
// @Tags kubernetes
// @Summary 获取 Pod 文件列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body pods.PodsFilesRequest true "获取 Pod 文件列表"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /kubernetes/pods/listFiles [post]
func (k *K8sPodApi) ListPodFiles(c *gin.Context) {
	req := pods.PodsFilesRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if files, err := k8sPodService.ListPodFiles(req); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(pods.PodFilesResponse{Files: files}, c)
	}
}

// UploadFiles 上传文件到 Pod
// @Tags kubernetes
// @Summary 上传文件到 Pod
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce application/json
// @Param data body pods.PodsFilesRequest true "上传文件到 Pod"
// @Success 200 {object} response.Response{msg=string} "上传成功"
// @Router /kubernetes/pods/uploadFile [post]
func (k *K8sPodApi) UploadFiles(c *gin.Context) {
	req := pods.PodsFilesRequest{}
	_ = c.ShouldBindQuery(&req)
	if err := utils.Verify(req, utils.PodVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	srcPath := filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().UnixNano()))
	err := saveTarFile(c, srcPath)
	if err != nil {
		global.DYCLOUD_LOG.Error("上传失败!", zap.Error(err))
		response.FailWithMessage("上传失败"+err.Error(), c)
		return
	}

	req.FilePath = srcPath
	err = k8sPodService.UploadFile(req)
	if err != nil {
		global.DYCLOUD_LOG.Error("上传失败!", zap.Error(err))
		response.FailWithMessage("上传失败"+err.Error(), c)
	} else {
		response.OkWithMessage("上传成功", c)
	}

}

// DeleteFiles 删除 Pod 文件
// @Tags kubernetes
// @Summary 删除 Pod 文件
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body pods.PodsFilesRequest true "删除 Pod 文件"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /kubernetes/pods/deleteFiles [post]
func (k *K8sPodApi) DeleteFiles(c *gin.Context) {
	req := pods.PodsFilesRequest{}
	_ = c.ShouldBindJSON(&req)
	if err := utils.Verify(req, utils.PodVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := k8sPodService.DeleteFile(req)
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)

	}
}
func saveTarFile(c *gin.Context, srcPath string) error {
	form, err := c.MultipartForm()
	if err != nil {
		global.DYCLOUD_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
	}

	files := form.File["file"]
	fw, err := os.Create(srcPath)
	defer fw.Close()
	if err != nil {
		return err
	}

	tw := tar.NewWriter(fw)
	defer tw.Close()

	for _, f := range files {
		hdr := &tar.Header{
			Name: f.Filename,
			Mode: 0644,
			Size: f.Size,
		}
		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}
		_f, err := f.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, _f)
		if err != nil {
			return err
		}
	}

	return nil
}
