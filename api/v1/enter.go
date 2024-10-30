package v1

import (
	"DYCLOUD/api/v1/cicd"
	"DYCLOUD/api/v1/cmdb"
	"DYCLOUD/api/v1/example"
	"DYCLOUD/api/v1/kubernetes/cloudtty"
	"DYCLOUD/api/v1/kubernetes/clusterManager/cluster"
	"DYCLOUD/api/v1/kubernetes/configManager/configmap"
	"DYCLOUD/api/v1/kubernetes/configManager/horizontalpodautoscalers"
	"DYCLOUD/api/v1/kubernetes/configManager/limitRange"
	"DYCLOUD/api/v1/kubernetes/configManager/poddisruptionbudget"
	"DYCLOUD/api/v1/kubernetes/configManager/resourceQuota"
	"DYCLOUD/api/v1/kubernetes/configManager/secret"
	"DYCLOUD/api/v1/kubernetes/namespaceManager/namespaces"
	"DYCLOUD/api/v1/kubernetes/networks/endpoint"
	"DYCLOUD/api/v1/kubernetes/networks/ingress"
	"DYCLOUD/api/v1/kubernetes/networks/service"
	"DYCLOUD/api/v1/kubernetes/nodeManager/node"
	"DYCLOUD/api/v1/kubernetes/rolesMangager/clusterrole"
	"DYCLOUD/api/v1/kubernetes/rolesMangager/clusterrolebinding"
	"DYCLOUD/api/v1/kubernetes/rolesMangager/rolebindings"
	"DYCLOUD/api/v1/kubernetes/rolesMangager/roles"
	"DYCLOUD/api/v1/kubernetes/rolesMangager/serviceaccount"
	"DYCLOUD/api/v1/kubernetes/storageManager/pv"
	"DYCLOUD/api/v1/kubernetes/storageManager/pvc"
	"DYCLOUD/api/v1/kubernetes/storageManager/storageClass"
	"DYCLOUD/api/v1/kubernetes/workload/cronjob"
	"DYCLOUD/api/v1/kubernetes/workload/daemonSet"
	"DYCLOUD/api/v1/kubernetes/workload/deployment"
	"DYCLOUD/api/v1/kubernetes/workload/job"
	"DYCLOUD/api/v1/kubernetes/workload/pod"
	"DYCLOUD/api/v1/kubernetes/workload/replicaset"
	"DYCLOUD/api/v1/kubernetes/workload/statefulSet"
	"DYCLOUD/api/v1/system"
	"DYCLOUD/api/v1/ws"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup      system.ApiGroup
	ExampleApiGroup     example.ApiGroup
	CmdbApiGroup        cmdb.ApiGroup
	ClusterApiGroup     cluster.ApiGroup
	NodeApiGroup        node.ApiGroup
	PodApiGroup         pod.ApiGroup
	DeploymentGroup     deployment.ApiGroup
	ConfigMapGroup      configmap.ApiGroup
	PvcGroup            pvc.ApiGroup
	ServiceAccount      serviceaccount.ApiGroup
	Secret              secret.ApiGroup
	Ingress             ingress.ApiGroup
	Service             service.ApiGroup
	Replicaset          replicaset.ApiGroup
	WsApi               ws.ApiGroup
	DaemonSet           daemonSet.ApiGroup
	StatefulSet         statefulSet.ApiGroup
	Job                 job.ApiGroup
	CronJob             cronjob.ApiGroup
	CloudTTY            cloudtty.ApiGroup
	Namespace           namespaces.ApiGroup
	Endpoint            endpoint.ApiGroup
	ResourceQuota       resourceQuota.ApiGroup
	LimitRange          limitRange.ApiGroup
	HorizontalPod       horizontalpodautoscalers.ApiGroup
	Poddisruptionbudget poddisruptionbudget.ApiGroup
	Pv                  pv.ApiGroup
	StorageClass        storageClass.ApiGroup
	ClusterRole         clusterrole.ApiGroup
	ClusterRoleBinding  clusterrolebinding.ApiGroup
	Role                roles.ApiGroup
	RoleBinding         rolebindings.ApiGroup
	CICD                cicd.ApiGroup
}
