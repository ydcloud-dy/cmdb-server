package cmdb

import api "DYCLOUD/api/v1"

type RouterGroup struct {
	CmdbProjectsRouter
	CmdbHostsRouter
	BatchOperationsRouter
}

var cmdbProjectsApi = api.ApiGroupApp.CmdbApiGroup.CmdbProjectsApi
var cmdbHostsApi = api.ApiGroupApp.CmdbApiGroup.CmdbHostsApi
var batchOperationsApi = api.ApiGroupApp.CmdbApiGroup.BatchOperationsApi
