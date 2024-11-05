package configCenter

import (
	api "DYCLOUD/api/v1"
)

type RouterGroup struct {
	EnvironmentRouter
	ServiceIntegrationRouter
	SourceCodeRouter
}

var EnvironmentApi = api.ApiGroupApp.ConfigCenter.EnvironmentApi
var ServiceIntegrationApi = api.ApiGroupApp.ConfigCenter.ServiceIntegrationApi
var SourceCodeApi = api.ApiGroupApp.ConfigCenter.SourceCodeApi
