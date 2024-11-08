package cicd

import api "DYCLOUD/api/v1"

type RouterGroup struct {
	ApplicationsRouter
	PipelinesRouter
}

var ApplicationApi = api.ApiGroupApp.CICD.ApplicationsApi

var PipelineApi = api.ApiGroupApp.CICD.PipelinesApi
