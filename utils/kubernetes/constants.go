package kubernetes

import (
	"time"

	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//  Copyright (c) 2014-2023 FIT2CLOUD

const (
	LabelManageKey    = "devops/manage"
	LabelRoleTypeKey  = "devops/role-type"
	LabelClusterId    = "devops/cluster-id"
	LabelUsername     = "devops/username"
	RoleTypeCluster   = "cluster"
	RoleTypeNamespace = "namespace"
)

var InitClusterRoles = []rbacV1.ClusterRole{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster-owner",
			Annotations: map[string]string{
				"description": "cluster-owner",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeCluster,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster-viewer",
			Annotations: map[string]string{
				"description": "cluster-viewer",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-namespaces",
			Annotations: map[string]string{
				"description": "manage-namespaces",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-namespaces",
			Annotations: map[string]string{
				"description": "view-namespaces",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-nodes",
			Annotations: map[string]string{
				"description": "view-nodes",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"nodes"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-nodes",
			Annotations: map[string]string{
				"description": "manage-nodes",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"nodes"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-crd",
			Annotations: map[string]string{
				"description": "view-crd",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"apiextensions.k8s.io"},
				Resources: []string{"customresourcedefinitions"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-crd",
			Annotations: map[string]string{
				"description": "manage-crd",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"apiextensions.k8s.io"},
				Resources: []string{"customresourcedefinitions"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-events",
			Annotations: map[string]string{
				"description": "view-events",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"events"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-cluster-rbac",
			Annotations: map[string]string{
				"description": "manage-cluster-rbac",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{"clusterroles", "clusterrolebindings"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-cluster-rbac",
			Annotations: map[string]string{
				"description": "view-cluster-rbac",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{"clusterroles", "clusterrolebindings"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-cluster-storage",
			Annotations: map[string]string{
				"description": "manage-cluster-storage",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"persistentvolumes"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{"storageclasses"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-cluster-storage",
			Annotations: map[string]string{
				"description": "view-cluster-storage",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeCluster},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"persistentvolumes"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{"storageclasses"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-serviceaccount-discovery",
			Annotations: map[string]string{
				"description": "manage-serviceaccount-discovery",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"services", "endpoints"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingresses", "networkpolicies"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-serviceaccount-discovery",
			Annotations: map[string]string{
				"description": "view-serviceaccount-discovery",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeNamespace},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"services", "endpoints"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingresses", "networkpolicies"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-config",
			Annotations: map[string]string{
				"description": "manage-config",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeNamespace},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps", "secrets", "resourcequotas", "limitranges"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"autoscaling"},
				Resources: []string{"horizontalpodautoscalers"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"policy"},
				Resources: []string{"poddisruptionbudgets"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-config",
			Annotations: map[string]string{
				"description": "view-config",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeNamespace},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps", "secrets", "resourcequotas", "limitranges"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"autoscaling"},
				Resources: []string{"horizontalpodautoscalers"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"policy"},
				Resources: []string{"poddisruptionbudgets"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-storage",
			Annotations: map[string]string{
				"description": "i18n_manage_storage",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"persistentvolumeclaims"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-storage",
			Annotations: map[string]string{
				"description": "view-storage",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"persistentvolumeclaims"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-workload",
			Annotations: map[string]string{
				"description": "manage-workload",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "pods/exec", "pods/log", "containers"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"deployments", "daemonsets", "replicasets", "statefulsets"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"batch"},
				Resources: []string{"jobs", "cronjobs"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-workload",
			Annotations: map[string]string{
				"description": "view-workload",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "containers", "pods/log"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"deployments", "daemonsets", "replicasets", "statefulsets"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"batch"},
				Resources: []string{"jobs", "cronjobs"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "view-rbac",
			Annotations: map[string]string{
				"description": "view-rbac",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{"roles", "rolebindings"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"serviceaccounts"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"policy"},
				Resources: []string{"podsecuritypolicies"},
				Verbs:     []string{"list", "get", "watch"},
			},
			{
				APIGroups: []string{"policy"},
				Resources: []string{"podsecuritypolicies"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-rbac",
			Annotations: map[string]string{
				"description": "manage-rbac",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{"roles", "rolebindings"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"serviceaccounts"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"policy"},
				Resources: []string{"podsecuritypolicies"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "namespace-owner",
			Annotations: map[string]string{
				"description": "namespace-owner",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "namespace-viewer",
			Annotations: map[string]string{
				"description": "namespace-viewer",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeNamespace,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "manage-appmarket",
			Annotations: map[string]string{
				"description": "manage-appmarket",
				"builtin":     "true",
				"created-at":  time.Now().Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				LabelManageKey:   "devops",
				LabelRoleTypeKey: RoleTypeCluster,
			},
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{"devops"},
				Resources: []string{"appmarkets"},
				Verbs:     []string{"*"},
			},
		},
	},
}
