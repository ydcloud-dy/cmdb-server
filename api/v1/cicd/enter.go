package cicd

import "DYCLOUD/service"

type ApiGroup struct {
	ApplicationsApi
	PipelinesApi
}

var ApplicationService = service.ServiceGroupApp.CICDServiceGroup.ApplicationsService
var PipelineService = service.ServiceGroupApp.CICDServiceGroup.PipelinesService
