package response

import (
	"DYCLOUD/model/kubernetes"
	cluster2 "DYCLOUD/model/kubernetes/cluster"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/rbac/v1"
)

type ClusterResponse struct {
	Cluster cluster2.K8sCluster `json:"cluster"`
}

type ClusterUserResponse struct {
	User []cluster2.User `json:"user"`
}

type RolesResponse struct {
	Roles []v1.ClusterRole `json:"roles"`
}

type ApiGroupResponse struct {
	Groups []kubernetes.ApiGroupOption `json:"groups"`
}

type ClusterUserNamespace struct {
	Namespaces []string `json:"namespaces"`
}

type ClusterListNamespace struct {
	Namespaces []corev1.Namespace `json:"namespaces"`
}
