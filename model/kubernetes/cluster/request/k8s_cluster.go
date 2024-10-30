package request

import (
	"DYCLOUD/model/common/request"
	"DYCLOUD/model/kubernetes/cluster"
	"time"
)

type K8sClusterSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	Alias          string     `json:"alias" form:"alias" `
	City           string     `json:"city" form:"city" `
	District       string     `json:"district" form:"district" `
	Name           string     `json:"name" form:"name" `
	request.PageInfo
}
type ClusterRoleType struct {
	RoleType  string `json:"role_type"`
	ClusterId uint   `json:"cluster_id"`
}

type ClusterApiGroups struct {
	ApiType   string `json:"api_type"`
	ClusterId uint   `json:"cluster_id"`
}

type NamespaceRoles struct {
	Namespace string   `json:"namespace"`
	Roles     []string `json:"roles"`
}

// CreateClusterRole
// @Description: 创建集群角色
type CreateClusterRole struct {
	ClusterRoles   []string         `json:"cluster_roles"`
	NamespaceRoles []NamespaceRoles `json:"namespace_roles"`
	UserUuids      []string         `json:"user_uuids"`
	cluster.User
	ClusterId int `json:"cluster_id"`
}

type DeleteClusterRole struct {
	UserUuids []string `json:"user_uuids"`
	ClusterId int      `json:"cluster_id"`
}
