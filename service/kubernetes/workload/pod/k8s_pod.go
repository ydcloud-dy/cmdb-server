package pod

import (
	"DYCLOUD/global"
	cluster2 "DYCLOUD/model/kubernetes/cluster"
	"DYCLOUD/model/kubernetes/pods"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"DYCLOUD/utils/kubernetes/podtool"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type K8sPodService struct{}

// GetPodList 获取 Pod 列表
//
// @Description 根据集群id和namespace获取pod列表
// @Tags Pod
// @Accept json
// @Produce json
// @Param req body pods.PodListReq true "请求参数"
// @Success 200 {object} []corev1.Pod
// @Router /pod/list [get]
func (k *K8sPodService) GetPodList(req pods.PodListReq) (podList *[]corev1.Pod, total int, err error) {
	var clusterIns = cluster2.K8sCluster{}
	if err := global.DYCLOUD_DB.Where(req.ClusterId).First(&clusterIns).Error; err != nil {
		return &[]corev1.Pod{}, 0, err
	}
	kubernetes := kubernetes.NewKubernetes(&clusterIns, &cluster2.User{}, true)
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("获取client-go客户端失败" + err.Error())
		return &[]corev1.Pod{}, 0, err
	}
	listOptions := metav1.ListOptions{FieldSelector: req.FieldSelector, LabelSelector: req.LabelSelector}
	data, err := client.CoreV1().Pods(req.Namespace).List(context.TODO(), listOptions)
	if err != nil {
		return &[]corev1.Pod{}, 0, err
	}
	var filterPod []corev1.Pod
	if req.Keyword != "" {
		for _, PV := range data.Items {
			if strings.Contains(PV.Name, req.Keyword) {
				filterPod = append(filterPod, PV)
			}
		}
	} else {
		filterPod = data.Items
	}
	if req.Page == 0 || req.PageSize == 0 {

		return &filterPod, len(filterPod), nil
	}
	result, err := paginate.Paginate(filterPod, req.Page, req.PageSize)

	fmt.Println(result)
	return result.(*[]corev1.Pod), len(filterPod), nil
}

// GetPodMetricsList 获取 Pod 监控指标列表
//
// @Description 根据集群id和namespace获取pod的监控指标列表
// @Tags Pod
// @Accept json
// @Produce json
// @Param req body pods.PodMetricsReq true "请求参数"
// @Success 200 {object} v1beta1.PodMetricsList
// @Router /pod/metrics [get]
func (k *K8sPodService) GetPodMetricsList(req pods.PodMetricsReq) (*v1beta1.PodMetricsList, error) {
	// 先到库里根据clusterID查询到集群实例
	var clusterIns = &cluster2.K8sCluster{}
	if err := global.DYCLOUD_DB.Where("id = ?", req.ClusterId).First(clusterIns).Error; err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &v1beta1.PodMetricsList{}, err
	}
	// 根据集群实例生成MetricsClientset对象
	kubernetes := kubernetes.NewNodeMetrics(clusterIns, &cluster2.User{}, true)
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &v1beta1.PodMetricsList{}, err
	}
	// 查询pod Metrics列表并返回
	podMetrics, err := client.MetricsV1beta1().PodMetricses(req.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &v1beta1.PodMetricsList{}, nil
	}
	return podMetrics, nil
}

// DescribePodInfo 查看 Pod 详情
//
// @Description 根据集群id和namespace获取pod详细信息
// @Tags Pod
// @Accept json
// @Produce json
// @Param req body pods.DescribePodInfo true "请求参数"
// @Success 200 {object} corev1.Pod
// @Router /pod/describe [get]
func (k *K8sPodService) DescribePodInfo(req pods.DescribePodInfo) (*corev1.Pod, error) {
	// 先到库里根据clusterID查询到集群实例
	var clusterIns = &cluster2.K8sCluster{}
	if err := global.DYCLOUD_DB.Where("id = ?", req.ClusterId).First(clusterIns).Error; err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Pod{}, err
	}
	// 根据集群实例生成Clientset对象
	kubernetes := kubernetes.NewKubernetes(clusterIns, &cluster2.User{}, true)
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Pod{}, err
	}
	// 查询pod 详情并返回
	podInfo, err := client.CoreV1().Pods(req.Namespace).Get(context.TODO(), req.PodName, metav1.GetOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Pod{}, err
	}
	return podInfo, nil
}

// PodEvents 获取 Pod 事件列表
//
// @Description 根据集群id和namespace获取pod的事件列表
// @Tags Pod
// @Accept json
// @Produce json
// @Param req body pods.PodEventsReq true "请求参数"
// @Success 200 {object} []corev1.Event
// @Router /pod/events [get]
func (k *K8sPodService) PodEvents(req pods.PodEventsReq) (*[]corev1.Event, int, error) {
	// 先到库里根据clusterID查询到集群实例
	var clusterIns = &cluster2.K8sCluster{}
	if err := global.DYCLOUD_DB.Where("id = ?", req.ClusterId).First(clusterIns).Error; err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &[]corev1.Event{}, 0, err
	}
	// 根据集群实例生成Clientset对象
	kubernetes := kubernetes.NewKubernetes(clusterIns, &cluster2.User{}, true)
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &[]corev1.Event{}, 0, err
	}
	// 查询特定 Pod 的事件信息
	options := metav1.ListOptions{
		FieldSelector: req.FieldSelector,
	}
	// 查询pod 详情并返回
	podInfo, err := client.CoreV1().Events(req.Namespace).List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &[]corev1.Event{}, 0, err
	}
	return &podInfo.Items, len(podInfo.Items), nil
}

// CreatePod 创建 Pod
//
// @Description 根据请求参数创建一个新的pod
// @Tags Pod
// @Accept json
// @Produce json
// @Param req body pods.CreatePodReq true "请求参数"
// @Success 200 {object} corev1.Pod
// @Router /pod/create [post]
func (p *K8sPodService) CreatePod(req pods.CreatePodReq) (*corev1.Pod, error) {
	var clusterIns = &cluster2.K8sCluster{}
	if err := global.DYCLOUD_DB.Where("id = ?", req.ClusterId).First(clusterIns).Error; err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Pod{}, err
	}
	// 根据集群实例生成Clientset对象
	kubernetes := kubernetes.NewKubernetes(clusterIns, &cluster2.User{}, true)
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Pod{}, err
	}
	tmp, err := json.Marshal(&req.Content)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Pod{}, err
	}
	data := string(tmp)
	var pod *corev1.Pod
	json.Unmarshal([]byte(data), &pod)
	result, err := client.CoreV1().Pods(req.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Pod{}, err
	}
	return result, nil
}

// DeletePod 删除 Pod
//
// @Description 根据请求参数删除一个pod
// @Tags Pod
// @Accept json
// @Produce json
// @Param req body pods.DeletePodReq true "请求参数"
// @Success 200 {string} string "删除成功"
// @Router /pod/delete [post]
func (p *K8sPodService) DeletePod(req pods.DeletePodReq) error {
	var clusterIns = &cluster2.K8sCluster{}
	if err := global.DYCLOUD_DB.Where("id = ?", req.ClusterId).First(clusterIns).Error; err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}
	// 根据集群实例生成Clientset对象
	kubernetes := kubernetes.NewKubernetes(clusterIns, &cluster2.User{}, true)
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}
	err = client.CoreV1().Pods(req.Namespace).Delete(context.TODO(), req.PodName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}
	return nil
}

// UpdatePod 更新 Pod
//
// @Description 根据请求参数更新一个pod
// @Tags Pod
// @Accept json
// @Produce json
// @Param req body pods.UpdatePodReq true "请求参数"
// @Success 200 {object} corev1.Pod
// @Router /pod/update [put]
func (p *K8sPodService) UpdatePod(req pods.UpdatePodReq) (*corev1.Pod, error) {
	var clusterIns = &cluster2.K8sCluster{}
	if err := global.DYCLOUD_DB.Where("id = ?", req.ClusterId).First(clusterIns).Error; err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Pod{}, err
	}
	// 根据集群实例生成Clientset对象
	kubernetes := kubernetes.NewKubernetes(clusterIns, &cluster2.User{}, true)
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Pod{}, err
	}
	tmp, err := json.Marshal(&req.Content)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Pod{}, err
	}
	data := string(tmp)
	var pod *corev1.Pod
	json.Unmarshal([]byte(data), &pod)
	result, err := client.CoreV1().Pods(req.Namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &corev1.Pod{}, err
	}
	return result, nil
}

// ListPodFiles 列出 Pod 中的文件
//
// @Description 列出 pod 中的文件
// @Tags Pod
// @Accept json
// @Produce json
// @Param req body pods.PodsFilesRequest true "请求参数"
// @Success 200 {object} []podtool.File
// @Router /pod/files/list [get]
func (p *K8sPodService) ListPodFiles(req pods.PodsFilesRequest) (files []podtool.File, err error) {
	pt, err := p.GetPodTool(req)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}
	return pt.ListFiles(req.Path)
}

// DownloadFile 从 Pod 中下载文件
//
// @Description 从 pod 中下载文件
// @Tags Pod
// @Accept json
// @Produce json
// @Param req body pods.PodsFilesRequest true "请求参数"
// @Success 200 {string} string "文件路径"
// @Router /pod/files/download [post]
func (p *K8sPodService) DownloadFile(req pods.PodsFilesRequest) (file string, err error) {
	var fileP string
	pt, err := p.GetPodTool(req)
	if err != nil {
		return fileP, err
	}
	fileNameWithSuffix := path.Base(req.Path)
	fileType := path.Ext(fileNameWithSuffix)
	fileName := strings.TrimSuffix(fileNameWithSuffix, fileType)
	fileP = filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().UnixNano()))
	err = os.MkdirAll(fileP, os.ModePerm)
	if err != nil {
		return "", err
	}
	fileP = filepath.Join(fileP, fileName+".tar")
	err = pt.CopyFromPod(req.Path, fileP)
	if err != nil {
		return "", err
	}

	return fileP, nil
}

// GetPodTool 获取 Pod 工具
//
// @Description 获取 pod 工具
// @Tags Pod
// @Accept json
// @Produce json
// @Param req body pods.PodsFilesRequest true "请求参数"
// @Success 200 {object} podtool.PodTool
// @Router /pod/tool [get]
func (p *K8sPodService) GetPodTool(req pods.PodsFilesRequest) (pt podtool.PodTool, err error) {

	var cl cluster2.K8sCluster
	if err := global.DYCLOUD_DB.Where("id = ?", req.ClusterId).First(&cl).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return pt, err
	}

	k := kubernetes.NewKubernetes(&cl, &cluster2.User{}, true)
	config, err := k.Config()
	if err != nil {
		return pt, err
	}

	clientSet, err := k.Client()
	if err != nil {
		return pt, err
	}

	pt = podtool.PodTool{
		Namespace:     req.Namespace,
		PodName:       req.PodName,
		ContainerName: req.ContainerName,
		K8sClient:     clientSet,
		RestClient:    config,
		ExecConfig: podtool.ExecConfig{
			Stdin: req.Stdin,
		},
	}

	return pt, nil
}

// UploadFile 上传文件到 Pod
//
// @Description 上传文件到 pod
// @Tags Pod
// @Accept json
// @Produce json
// @Param req body pods.PodsFilesRequest true "请求参数"
// @Success 200 {string} string "上传成功"
// @Router /pod/files/upload [post]
func (k *K8sPodService) UploadFile(req pods.PodsFilesRequest) error {
	pt, err := k.GetPodTool(req)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}
	reader, writer := io.Pipe()
	pt.ExecConfig.Stdin = reader
	go func() {
		defer func() {
			_ = writer.Close()
		}()
		tarFile, err := os.Open(req.FilePath)
		if err != nil {
			global.DYCLOUD_LOG.Error(err.Error())
			return
		}
		_, err = io.Copy(writer, tarFile)
		if err != nil {
			global.DYCLOUD_LOG.Error(err.Error())
			return
		}
	}()
	return pt.CopyToContainer(req.Path)

}

// DeleteFile 删除 Pod 中的文件
//
// @Description 删除 pod 中的文件
// @Tags Pod
// @Accept json
// @Produce json
// @Param req body pods.PodsFilesRequest true "请求参数"
// @Success 200 {string} string "删除成功"
// @Router /pod/files/delete [post]
func (p *K8sPodService) DeleteFile(req pods.PodsFilesRequest) (err error) {
	req.Commands = []string{"rm", "-rf", req.Path}
	pt, err := p.GetPodTool(req)
	if err != nil {
		return err
	}

	_, err = pt.ExecCommand(req.Commands)
	if err != nil {
		return err
	}

	return
}
