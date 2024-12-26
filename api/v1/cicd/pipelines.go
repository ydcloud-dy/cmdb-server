package cicd

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cicd"
	request2 "DYCLOUD/model/cicd/request"
	"DYCLOUD/model/common/response"
	"DYCLOUD/service/kubernetes/cluster"
	"DYCLOUD/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"strconv"
)

type PipelinesApi struct{}

// GetPipelinesList
//
//	@Description: 获取应用列表
//	@receiver PipelinesApi
//	@param c
func (PipelinesApi *PipelinesApi) GetPipelinesList(c *gin.Context) {
	var env *request2.PipelinesRequest
	err := c.BindQuery(&env)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//Pipelines.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, total, err := PipelineService.GetPipelinesList(env)
	if err != nil {
		global.DYCLOUD_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     data,
		Total:    total,
		Page:     env.Page,
		PageSize: env.PageSize,
	}, "获取成功", c)
}

// DescribePipelines
//
//	@Description: 查看环境详情
//	@receiver PipelinesApi
//	@param c
func (PipelinesApi *PipelinesApi) DescribePipelines(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//Pipelines.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, err := PipelineService.DescribePipelines(id)
	if err != nil {
		global.DYCLOUD_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

func (PipelinesApi *PipelinesApi) GetPipelinesStatus(c *gin.Context) {
	var request *request2.PipelinesRequest
	err := c.BindQuery(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 从请求中获取集群名称
	clusterName := request.Cluster_ID
	if clusterName == 0 {
		response.FailWithMessage("集群ID不能为空", c)
		return
	}
	k8sService := cluster.K8sClusterService{}
	cluster, err := k8sService.GetK8sCluster(request.Cluster_ID)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取集群信息失败: %v", err), c)
		return
	}
	if cluster.KubeConfig == "" {
		response.FailWithMessage("集群的 kubeConfig 不能为空", c)
		return
	}
	// 解析 kubeConfig 内容
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.KubeConfig))
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("加载 kubeConfig 失败: %v", err), c)
		return
	}

	// 初始化 Tekton 客户端
	clientset, err := tektonclient.NewForConfig(config)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建 Tekton 客户端失败", err.Error()), c)
		return
	}
	data, err := PipelineService.GetPipelinesStatus(clientset, request)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取pipeline失败", err.Error()), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)

}

// CreatePipelines
//
//	@Description: 创建环境
//	@receiver PipelinesApi
//	@param c
func (PipelinesApi *PipelinesApi) CreatePipelines(c *gin.Context) {

	var request *cicd.Pipelines
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 从请求中获取集群名称
	clusterName := request.K8SClusterName
	if clusterName == "" {
		response.FailWithMessage("集群名称不能为空", c)
		return
	}
	k8sService := cluster.K8sClusterService{}
	cluster, err := k8sService.GetK8sClusterByName(request.K8SClusterName)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取集群信息失败: %v", err), c)
		return
	}
	if cluster.KubeConfig == "" {
		response.FailWithMessage("集群的 kubeConfig 不能为空", c)
		return
	}
	// 解析 kubeConfig 内容
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.KubeConfig))
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("加载 kubeConfig 失败: %v", err), c)
		return
	}

	// 初始化 Tekton 客户端
	clientset, err := tektonclient.NewForConfig(config)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建 Tekton 客户端失败", err.Error()), c)
		return
	}
	request.CreatedBy = utils.GetUserID(c)
	request.CreatedName = utils.GetUserName(c)

	k8sClient, err := kubernetes.NewForConfig(config)

	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建 k8s 客户端失败", err.Error()), c)
		return
	}
	err = PipelineService.CreatePipelines(k8sClient, clientset, request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("创建成功", c)
}
func (PipelinesApi *PipelinesApi) RunPipelines(c *gin.Context) {

	var request *cicd.Pipelines
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 从请求中获取集群名称
	clusterName := request.K8SClusterName
	if clusterName == "" {
		response.FailWithMessage("集群名称不能为空", c)
		return
	}
	k8sService := cluster.K8sClusterService{}
	cluster, err := k8sService.GetK8sClusterByName(request.K8SClusterName)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取集群信息失败: %v", err), c)
		return
	}
	if cluster.KubeConfig == "" {
		response.FailWithMessage("集群的 kubeConfig 不能为空", c)
		return
	}
	// 解析 kubeConfig 内容
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.KubeConfig))
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("加载 kubeConfig 失败: %v", err), c)
		return
	}

	// 初始化 Tekton 客户端
	clientset, err := tektonclient.NewForConfig(config)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建 Tekton 客户端失败", err.Error()), c)
		return
	}
	request.CreatedBy = utils.GetUserID(c)
	request.CreatedName = utils.GetUserName(c)

	k8sClient, err := kubernetes.NewForConfig(config)

	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建 k8s 客户端失败", err.Error()), c)
		return
	}
	err = PipelineService.RunPipelines(k8sClient, clientset, request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("创建成功", c)
}

// UpdatePipelines
//
//	@Description: 更新环境
//	@receiver PipelinesApi
//	@param c
func (PipelinesApi *PipelinesApi) UpdatePipelines(c *gin.Context) {
	var request *cicd.Pipelines
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 从请求中获取集群名称
	clusterName := request.K8SClusterName
	if clusterName == "" {
		response.FailWithMessage("集群名称不能为空", c)
		return
	}
	k8sService := cluster.K8sClusterService{}
	cluster, err := k8sService.GetK8sClusterByName(request.K8SClusterName)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取集群信息失败: %v", err), c)
		return
	}
	if cluster.KubeConfig == "" {
		response.FailWithMessage("集群的 kubeConfig 不能为空", c)
		return
	}
	// 解析 kubeConfig 内容
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.KubeConfig))
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("加载 kubeConfig 失败: %v", err), c)
		return
	}

	// 初始化 Tekton 客户端
	clientset, err := tektonclient.NewForConfig(config)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建 Tekton 客户端失败", err.Error()), c)
		return
	}
	request.CreatedBy = utils.GetUserID(c)
	request.CreatedName = utils.GetUserName(c)

	k8sClient, err := kubernetes.NewForConfig(config)

	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建 k8s 客户端失败", err.Error()), c)
		return
	}

	data, err := PipelineService.UpdatePipelines(k8sClient, clientset, request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

// DeletePipelines
//
//	@Description: 删除环境
//	@receiver PipelinesApi
//	@param c
func (PipelinesApi *PipelinesApi) DeletePipelines(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = PipelineService.DeletePipelines(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("删除成功", c)
}

// DeletePipelinesByIds
//
//	@Description: 根据id批量删除环境
//	@receiver PipelinesApi
//	@param c
func (PipelinesApi *PipelinesApi) DeletePipelinesByIds(c *gin.Context) {
	ids := &request2.DeleteApplicationByIds{}
	err := c.ShouldBindJSON(ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = PipelineService.DeletePipelinesByIds(ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("删除成功", c)
}

func (PipelinesApi *PipelinesApi) GetPipelinesNotice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	notice, err := PipelineService.GetPipelinesNotice(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(notice, c)

}
func (PipelinesApi *PipelinesApi) ClosePipelineNotice(c *gin.Context) {
	req := &request2.ClosePipelineNotice{}
	pipelineID, _ := strconv.Atoi(c.Param("id"))
	err := c.ShouldBindJSON(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = PipelineService.ClosePipelineNotice(req, pipelineID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("关闭通知成功", c)

}
func (PipelinesApi *PipelinesApi) CreatePipelineNotice(c *gin.Context) {
	req := &cicd.Notice{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = PipelineService.CreatePipelinesNotice(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("创建通知成功", c)

}

func (PipelinesApi *PipelinesApi) GetPipelinesCache(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	notice, err := PipelineService.GetPipelinesCache(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(notice, c)

}
func (PipelinesApi *PipelinesApi) ClosePipelineCache(c *gin.Context) {
	req := &request2.ClosePipelineCache{}
	pipelineID, _ := strconv.Atoi(c.Param("id"))
	err := c.ShouldBindJSON(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = PipelineService.ClosePipelineCache(req, pipelineID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("关闭缓存成功", c)

}
func (PipelinesApi *PipelinesApi) CreatePipelineCache(c *gin.Context) {
	req := &cicd.Cache{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = PipelineService.CreatePipelinesCache(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("创建通知成功", c)

}

func (PipelinesApi *PipelinesApi) SyncBranches(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//Applications.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	err = PipelineService.SyncBranches(id)
	if err != nil {
		global.DYCLOUD_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithData("同步成功", c)
}
func (PipelinesApi *PipelinesApi) GetBranchesList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var branches = request2.ApplicationRequest{}
	err = c.BindQuery(&branches)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	branches.AppId = id

	//Applications.CreatedBy = utils.GetUserID(c)
	//userId := utils.GetUserID(c)
	data, total, err := PipelineService.GetBranchesList(branches)
	if err != nil {
		global.DYCLOUD_LOG.Error("执行失败!", zap.Error(err))
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     data,
		Total:    total,
		Page:     branches.Page,
		PageSize: branches.PageSize,
	}, "获取成功", c)
}
