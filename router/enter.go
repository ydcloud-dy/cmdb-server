package router

import (
	"DYCLOUD/router/cmdb"
	"DYCLOUD/router/configCenter"
	"DYCLOUD/router/example"
	"DYCLOUD/router/kubernetes/cloudtty"
	"DYCLOUD/router/kubernetes/clusterManager/cluster"
	"DYCLOUD/router/kubernetes/configManager/configmap"
	"DYCLOUD/router/kubernetes/configManager/horizontalPod"
	"DYCLOUD/router/kubernetes/configManager/limitRange"
	"DYCLOUD/router/kubernetes/configManager/poddistruptionbudget"
	"DYCLOUD/router/kubernetes/configManager/resourceQuota"
	"DYCLOUD/router/kubernetes/configManager/secret"
	"DYCLOUD/router/kubernetes/namespaceManager/namespaces"
	"DYCLOUD/router/kubernetes/network/Ingress"
	"DYCLOUD/router/kubernetes/network/endpoint"
	"DYCLOUD/router/kubernetes/network/service"
	"DYCLOUD/router/kubernetes/nodeManager/node"
	"DYCLOUD/router/kubernetes/rolesManager/clusterRole"
	"DYCLOUD/router/kubernetes/rolesManager/clusterRoleBinding"
	"DYCLOUD/router/kubernetes/rolesManager/rolebinding"
	"DYCLOUD/router/kubernetes/rolesManager/roles"
	"DYCLOUD/router/kubernetes/rolesManager/serviceAccount"
	"DYCLOUD/router/kubernetes/storageManager/pv"
	"DYCLOUD/router/kubernetes/storageManager/pvc"
	"DYCLOUD/router/kubernetes/storageManager/storageClass"
	"DYCLOUD/router/kubernetes/workload/cronjob"
	"DYCLOUD/router/kubernetes/workload/daemonset"
	"DYCLOUD/router/kubernetes/workload/deployment"
	"DYCLOUD/router/kubernetes/workload/job"
	"DYCLOUD/router/kubernetes/workload/pod"
	"DYCLOUD/router/kubernetes/workload/replicaSet"
	"DYCLOUD/router/kubernetes/workload/statefulSet"
	"DYCLOUD/router/system"
	"DYCLOUD/router/velero"
	"DYCLOUD/router/ws"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System               system.RouterGroup
	Example              example.RouterGroup
	Cmdb                 cmdb.RouterGroup
	Cluster              cluster.RouterGroup
	Node                 node.RouterGroup
	Pod                  pod.RouterGroup
	Deployment           deployment.RouterGroup
	Secret               secret.RouterGroup
	Configmap            configmap.RouterGroup
	ServiceAccount       serviceAccount.RouterGroup
	Pvc                  pvc.RouterGroup
	Service              service.RouterGroup
	Ingress              Ingress.RouterGroup
	ReplicaSet           replicaSet.RouterGroup
	Ws                   ws.RouterGroup
	DaemonSet            daemonset.RouterGroup
	StatefulSet          statefulSet.RouterGroup
	Job                  job.RouterGroup
	CronJob              cronjob.RouterGroup
	CloudTTY             cloudtty.RouterGroup
	Namespace            namespaces.RouterGroup
	EndPoint             endpoint.RouterGroup
	ResourceQuota        resourceQuota.RouterGroup
	LimitRange           limitRange.RouterGroup
	HorizontalPod        horizontalPod.RouterGroup
	Poddistruptionbudget poddistruptionbudget.RouterGroup
	Pv                   pv.RouterGroup
	StorageClass         storageClass.RouterGroup
	ClusterRole          clusterRole.RouterGroup
	ClusterRoleBinding   clusterRoleBinding.RouterGroup
	Role                 roles.RouterGroup
	RoleBinding          rolebinding.RouterGroup
	Velero               velero.K8sVeleroTasksRouter
	Environment          configCenter.RouterGroup
}
