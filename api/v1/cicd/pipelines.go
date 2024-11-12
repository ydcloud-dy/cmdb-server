package cicd

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cicd"
	request2 "DYCLOUD/model/cicd/request"
	"DYCLOUD/model/common/response"
	"fmt"
	"github.com/gin-gonic/gin"
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
	"strconv"
)

type PipelinesApi struct{}

// GetPipelinesList
//
//	@Description: 获取应用列表
//	@receiver PipelinesApi
//	@param c
func (PipelinesApi *PipelinesApi) GetPipelinesList(c *gin.Context) {
	var env *request2.ApplicationRequest
	err := c.BindQuery(&env)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(env)
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
	fmt.Println(id)
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
	fmt.Println(request)
	//request.CreatedBy = utils.GetUserID(c)
	//request.CreatedName = utils.GetUserName(c)
	//userId := utils.GetUserID(c)
	configPath := filepath.Join("/Users/dujie/.kube", "config") // 替换为你测试时的实际路径
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("加载 kubeconfig 失败", err.Error()), c)
		return
	}

	// 初始化 Tekton 客户端
	clientset, err := tektonclient.NewForConfig(config)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建 Tekton 客户端失败", err.Error()), c)
		return
	}
	err = PipelineService.CreatePipelines(clientset, request)
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
	fmt.Println(request, "================================")
	//request.UpdatedBy = utils.GetUserID(c)
	//request.UpdatedName = utils.GetUserName(c)
	data, err := PipelineService.UpdatePipelines(request)
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
	fmt.Println(id)
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
	fmt.Println(ids)
	err = PipelineService.DeletePipelinesByIds(ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("删除成功", c)
}
