package configCenter

import (
	"DYCLOUD/service"
)

type ApiGroup struct {
	EnvironmentApi
	ServiceIntegrationApi
	SourceCodeApi
	BuildEnvApi
}

var EnvironmentService = service.ServiceGroupApp.ConfigCenterServiceGroup.EnvironmentService
var ServiceIntegrationService = service.ServiceGroupApp.ConfigCenterServiceGroup.ServiceIntegrationService
var SourceCodeService = service.ServiceGroupApp.ConfigCenterServiceGroup.SourceCodeService
var BuildEnvService = service.ServiceGroupApp.ConfigCenterServiceGroup.BuildEnvService
