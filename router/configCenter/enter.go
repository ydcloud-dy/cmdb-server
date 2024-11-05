package configCenter

import api "DYCLOUD/api/v1"

type RouterGroup struct {
	EnvironmentRouter
	ServiceIntegrationRouter
	SourceCodeRouter
	ApplicationsRouter
}

var EnvironmentApi = api.ApiGroupApp.CICD.EnvironmentApi
var ServiceIntegrationApi = api.ApiGroupApp.CICD.ServiceIntegrationApi
var SourceCodeApi = api.ApiGroupApp.CICD.SourceCodeApi
var ApplicationApi = api.ApiGroupApp.CICD.ApplicationsApi
