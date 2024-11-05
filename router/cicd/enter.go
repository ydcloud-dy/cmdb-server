package cicd

import api "DYCLOUD/api/v1"

type RouterGroup struct {
	ApplicationsRouter
}

var ApplicationApi = api.ApiGroupApp.CICD.ApplicationsApi
