package service

import (
	"DYCLOUD/service/cicd"
	"DYCLOUD/service/cmdb"
	"DYCLOUD/service/configCenter"
	"DYCLOUD/service/example"
	"DYCLOUD/service/kubernetes/cloudtty"
	"DYCLOUD/service/kubernetes/cluster"
	"DYCLOUD/service/kubernetes/configManager/configmap"
	"DYCLOUD/service/kubernetes/configManager/horizontalPod"
	"DYCLOUD/service/kubernetes/configManager/limitRange"
	"DYCLOUD/service/kubernetes/configManager/poddistruptionbudget"
	"DYCLOUD/service/kubernetes/configManager/resourceQuota"
	"DYCLOUD/service/kubernetes/configManager/secret"
	"DYCLOUD/service/kubernetes/namespace"
	"DYCLOUD/service/kubernetes/network/Ingress"
	"DYCLOUD/service/kubernetes/network/endpoint"
	"DYCLOUD/service/kubernetes/network/service"
	"DYCLOUD/service/kubernetes/node"
	"DYCLOUD/service/kubernetes/rolesManager/clusterRole"
	"DYCLOUD/service/kubernetes/rolesManager/clusterRoleBinding"
	"DYCLOUD/service/kubernetes/rolesManager/role"
	"DYCLOUD/service/kubernetes/rolesManager/roleBinding"
	"DYCLOUD/service/kubernetes/serviceAccount"
	"DYCLOUD/service/kubernetes/storageManager/pv"
	"DYCLOUD/service/kubernetes/storageManager/pvc"
	"DYCLOUD/service/kubernetes/storageManager/storageClass"
	"DYCLOUD/service/kubernetes/velero"
	"DYCLOUD/service/kubernetes/workload/cronjob"
	"DYCLOUD/service/kubernetes/workload/daemonset"
	"DYCLOUD/service/kubernetes/workload/deployment"
	"DYCLOUD/service/kubernetes/workload/job"
	"DYCLOUD/service/kubernetes/workload/pod"
	"DYCLOUD/service/kubernetes/workload/replicaSet"
	"DYCLOUD/service/kubernetes/workload/statefulSet"
	"DYCLOUD/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup               system.ServiceGroup
	ExampleServiceGroup              example.ServiceGroup
	CmdbServiceGroup                 cmdb.ServiceGroup
	ClusterServiceGroup              cluster.ServiceGroup
	NodeServiceGroup                 node.ServiceGroup
	PodServiceGroup                  pod.ServiceGroup
	DeploymentServiceGroup           deployment.ServiceGroup
	PvcServiceGroup                  pvc.ServiceGroup
	SecretServiceGroup               secret.ServiceGroup
	ConfigMapServiceGroup            configmap.ServiceGroup
	ServiceAccountServiceGroup       serviceAccount.ServiceGroup
	IngressServiceGroup              Ingress.ServiceGroup
	SvcServiceGroup                  service.ServiceGroup
	ReplicaSetServiceGroup           replicaSet.ServiceGroup
	DaemonSetServiceGroup            daemonset.ServiceGroup
	StatefulSetServiceGroup          statefulSet.ServiceGroup
	JobServiceGroup                  job.ServiceGroup
	CronJobServiceGroup              cronjob.ServiceGroup
	CloudTTYServiceGroup             cloudtty.ServiceGroup
	NamespaceServiceGroup            namespace.ServiceGroup
	EndPointServiceGroup             endpoint.ServiceGroup
	ResourceQuotaServiceGroup        resourceQuota.ServiceGroup
	LimitRangeServiceGroup           limitRange.ServiceGroup
	HorizontalPodServiceGroup        horizontalPod.ServiceGroup
	PoddistruptionbudgetServiceGroup poddistruptionbudget.ServiceGroup
	PvServiceGroup                   pv.ServiceGroup
	StorageClassServiceGroup         storageClass.ServiceGroup
	ClusterRoleServiceGroup          clusterRole.ServiceGroup
	ClusterRoleBindingServiceGroup   clusterRoleBinding.ServiceGroup
	RoleServiceGroup                 role.ServiceGroup
	RoleBindingServiceGroup          roleBinding.ServiceGroup
	VeleroServiceGroup               velero.ServiceGroup
	ConfigCenterServiceGroup         configCenter.ServiceGroup
	CICDServiceGroup                 cicd.ServiceGroup
}
