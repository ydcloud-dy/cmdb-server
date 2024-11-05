package configCenter

import (
	"DYCLOUD/service"
)

type ApiGroup struct {
	EnvironmentApi
	ServiceIntegrationApi
	SourceCodeApi
}

var EnvironmentService = service.ServiceGroupApp.CICDServiceGroup.EnvironmentService
var ServiceIntegrationService = service.ServiceGroupApp.CICDServiceGroup.ServiceIntegrationService
var SourceCodeService = service.ServiceGroupApp.CICDServiceGroup.SourceCodeService
