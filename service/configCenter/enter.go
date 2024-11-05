package configCenter

import "DYCLOUD/service/cicd"

type ServiceGroup struct {
	EnvironmentService
	ServiceIntegrationService
	SourceCodeService
	cicd.ApplicationsService
}
