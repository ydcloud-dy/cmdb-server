package metrics

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/response"
	"DYCLOUD/model/kubernetes/metrics"
	"DYCLOUD/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MetricsApi struct{}

var metricsService = service.ServiceGroupApp.MetricsServiceGroup.MetricsService

// @Tags MetricsApi
// @Summary  MetricsGet
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body kubernetesReq.ResourceParamRequest true "分页获取API列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "分页获取API列表,返回包括列表,总数,页码,每页数量"
// @Router  /kubernetes/metrics/get [post]
func (m *MetricsApi) MetricsGet(c *gin.Context) {
	var metricsQuery metrics.MetricsQuery
	_ = c.ShouldBindJSON(&metricsQuery)

	if queryResp, err := metricsService.GetMetrics(metricsQuery); err != nil {
		global.DYCLOUD_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
	} else {
		response.OkWithDetailed(metrics.MetricsResponse{
			Metrics: queryResp,
		}, "获取成功", c)
	}
}
