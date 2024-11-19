package initialize

import (
	"DYCLOUD/router"
	"github.com/gin-gonic/gin"
)

func holder(routers ...*gin.RouterGroup) {
	_ = routers
	_ = router.RouterGroupApp
}
func initBizRouter(routers ...*gin.RouterGroup) {
	privateGroup := routers[0]
	publicGroup := routers[1]
	{
		clusterRouter := router.RouterGroupApp.Cluster
		clusterRouter.InitK8sClusterRouter(privateGroup, publicGroup)
		nodeRouter := router.RouterGroupApp.Node
		nodeRouter.Initk8sNodeRouter(privateGroup, publicGroup)
		podRouter := router.RouterGroupApp.Pod
		podRouter.Initk8sPodRouter(privateGroup, publicGroup)
		deploymentRouter := router.RouterGroupApp.Deployment
		deploymentRouter.InitK8sDeploymentRouter(privateGroup, publicGroup)
		secretRouter := router.RouterGroupApp.Secret
		secretRouter.Initk8sSecretRouter(privateGroup, publicGroup)
		saRouter := router.RouterGroupApp.ServiceAccount
		saRouter.Initk8sServiceAccountRouter(privateGroup, publicGroup)
		configMapRouter := router.RouterGroupApp.Configmap
		configMapRouter.InitK8sConfigMapRouter(privateGroup, publicGroup)
		pvcRouter := router.RouterGroupApp.Pvc
		pvcRouter.Initk8sPvcRouter(privateGroup, publicGroup)
		serviceRouter := router.RouterGroupApp.Service
		serviceRouter.Initk8sServiceRouter(privateGroup, publicGroup)
		ingressRouter := router.RouterGroupApp.Ingress
		ingressRouter.Initk8sIngressRouter(privateGroup, publicGroup)
		replicaSetRouter := router.RouterGroupApp.ReplicaSet
		replicaSetRouter.Initk8sReplicaSetRouter(privateGroup, publicGroup)
		ws := router.RouterGroupApp.Ws
		ws.InitWsRouter(publicGroup)
		daemonSetRouter := router.RouterGroupApp.DaemonSet
		daemonSetRouter.Initk8sDaemonSetRouter(privateGroup, publicGroup)
		statefulSetRouter := router.RouterGroupApp.StatefulSet
		statefulSetRouter.Initk8sStatefulSetRouter(privateGroup, publicGroup)
		jobRouter := router.RouterGroupApp.Job
		jobRouter.Initk8sJobRouter(privateGroup, publicGroup)
		cronJobRouter := router.RouterGroupApp.CronJob
		cronJobRouter.Initk8sCronJobRouter(privateGroup, publicGroup)
		cloudttyRouter := router.RouterGroupApp.CloudTTY
		cloudttyRouter.Initk8sCloudTTYRouter(privateGroup, publicGroup)
		cloudttyRouter.Initk8sNodeTTYRouter(privateGroup, publicGroup)
		namespaceRouter := router.RouterGroupApp.Namespace
		namespaceRouter.Initk8sNamespaceRouter(privateGroup, publicGroup)
		endpointRouter := router.RouterGroupApp.EndPoint
		endpointRouter.Initk8sEndPointRouter(privateGroup, publicGroup)
		resourceQuotaRouter := router.RouterGroupApp.ResourceQuota
		resourceQuotaRouter.Initk8sResourceQuotaRouter(privateGroup, publicGroup)
		horizontalPodRouter := router.RouterGroupApp.HorizontalPod
		horizontalPodRouter.Initk8sHorizontalPodRouter(privateGroup, publicGroup)
		PoddistruptionbudGetRouter := router.RouterGroupApp.Poddistruptionbudget
		PoddistruptionbudGetRouter.Initk8sPoddisruptionbudgetRouter(privateGroup, publicGroup)
		limitRangeRouter := router.RouterGroupApp.LimitRange
		limitRangeRouter.Initk8sLimitRangeRouter(privateGroup, publicGroup)
		pvRouter := router.RouterGroupApp.Pv
		pvRouter.Initk8sPVRouter(privateGroup, publicGroup)
		storageClassRouter := router.RouterGroupApp.StorageClass
		storageClassRouter.Initk8sStorageClassRouter(privateGroup, publicGroup)
		clusterRoleRouter := router.RouterGroupApp.ClusterRole
		clusterRoleRouter.Initk8sClusterRoleRouter(privateGroup, publicGroup)
		clusterRoleBindingRouter := router.RouterGroupApp.ClusterRoleBinding
		clusterRoleBindingRouter.Initk8sClusterRoleBindingRouter(privateGroup, publicGroup)
		roleRouter := router.RouterGroupApp.Role
		roleRouter.Initk8sRoleRouter(privateGroup, publicGroup)
		roleBindingRouter := router.RouterGroupApp.RoleBinding
		roleBindingRouter.Initk8sRoleBindingRouter(privateGroup, publicGroup)
		veleroRouter := router.RouterGroupApp.Velero
		veleroRouter.InitK8sVeleroTasksRouter(privateGroup, publicGroup)
	}
	holder(publicGroup, privateGroup) // 占位方法，保证文件可以正确加载，避免go空变量检测报错，请勿删除。
	{
		cmdbRouter := router.RouterGroupApp.Cmdb
		cmdbRouter.InitCmdbProjectsRouter(privateGroup, publicGroup)
		cmdbRouter.InitCmdbHostsRouter(privateGroup, publicGroup)
		cmdbRouter.InitBatchOperationsRouter(privateGroup, publicGroup)
	}
	{
		configCenter := router.RouterGroupApp.Environment
		configCenter.InitEnvironmentRouter(privateGroup, publicGroup)
		configCenter.InitServiceIntrgrationRouter(privateGroup, publicGroup)
		configCenter.InitSourceCodeRouter(privateGroup, publicGroup)
		configCenter.InitBuildEnvRouter(privateGroup, publicGroup)

	}
	{
		cicd := router.RouterGroupApp.CICD
		cicd.InitApplicationsRouter(privateGroup, publicGroup)
		cicd.InitPipelinesRouter(privateGroup, publicGroup)
	}
}
