// 自动生成模板K8sCluster
package cluster

import (
	"DYCLOUD/global"
	"github.com/gofrs/uuid/v5"
	rbacV1 "k8s.io/api/rbac/v1"
)

// k8sCluster表 结构体  K8sCluster
type K8sCluster struct {
	global.DYCLOUD_MODEL
	ID                 uint      `json:"id" gorm:"not null;unique;primary_key"`
	UUID               uuid.UUID `json:"uuid" gorm:"comment:集群UUID"`
	Name               string    `json:"name" form:"name" gorm:"comment:集群名称"`
	KubeType           uint      `json:"kube_type" form:"kube_type" gorm:"comment:凭据类型1:KubeConfig,2:Token"`
	KubeConfig         string    `gorm:"type:longText" json:"kube_config" form:"kube_config" gorm:"comment:kube_config"`
	KubeCaCrt          string    `gorm:"type:longText; comment:ca.crt" json:"kube_ca_crt" form:"kube_ca_crt"`
	ApiAddress         string    `gorm:"type:longText" json:"api_address" form:"api_address" gorm:"comment:api_address"`
	PrometheusUrl      string    `gorm:"type:longText" json:"prometheus_url" form:"prometheus_url" gorm:"comment:prometheus 地址"`
	PrometheusAuthType uint      `json:"prometheus_auth_type" form:"prometheus_auth_type" gorm:"comment: 认证类型"`
	PrometheusUser     string    `gorm:"type:longText" json:"prometheus_user" form:"prometheus_user" gorm:"comment:用户名"`
	PrometheusPwd      string    `gorm:"type:longText" json:"prometheus_pwd" form:"prometheus_pwd" gorm:"comment:密码"`
	Users              []User    `json:"users" gorm:"foreignKey:ClusterId;"`
	CreatedBy          uint      `gorm:"column:created_by;comment:创建者"`
	UpdatedBy          uint      `gorm:"column:updated_by;comment:更新者"`
	DeletedBy          uint      `gorm:"column:deleted_by;comment:删除者"`
}

type RoleData struct {
	ClusterId uint `json:"cluster_id"`
	rbacV1.ClusterRole
}
type GetClusterById struct {
	ID int `json:"id" form:"id"` // 主键ID
}

// TableName k8sCluster表 K8sCluster自定义表名 k8s_cluster
func (K8sCluster) TableName() string {
	return "k8s_clusters"
}
