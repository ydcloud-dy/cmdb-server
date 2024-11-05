package cicd

import "DYCLOUD/service"

type ApiGroup struct {
	ApplicationsApi
}

var ApplicationService = service.ServiceGroupApp.CICDServiceGroup.ApplicationsService
