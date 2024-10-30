package kubernetes

import (
	"DYCLOUD/global"
	cluster2 "DYCLOUD/model/kubernetes/cluster"
	"DYCLOUD/model/kubernetes/nodes"
	"DYCLOUD/model/system"
	"github.com/gofrs/uuid/v5"
	"reflect"
)

// KubernetesInterface
// @Description: 通用接口，返回当前集群实例、集群用户实例、是否是管理员
//
//	type KubernetesInterface interface {
//		Generic(req GenericInterface, uuid uuid.UUID) (*cluster.K8sCluster, *cluster.User, bool, error)
//	}
type ClusterIDGetter interface {
	GetClusterID() int
}
type BaseService struct {
	Kubernetes  *Kubernetes
	NodeMetrics *NodeMetrics
}

// Generic 通用操作
//
// @Description 通用操作，用于获取集群实例和用户实例
// @Tags Kubernetes
// @Accept json
// @Produce json
// @Param req body ClusterIDGetter true "请求参数"
// @Param uuid path string true "用户UUID"
// @Success 200 {object} Kubernetes
// @Router /kubernetes/generic [post]
func (d *BaseService) Generic(req ClusterIDGetter, uuid uuid.UUID) (*Kubernetes, error) {
	var clusterIns = &cluster2.K8sCluster{}
	k8sType := 1
	if clusterIDGetter, ok := req.(ClusterIDGetter); ok {
		clusterID := clusterIDGetter.GetClusterID()
		if err := global.DYCLOUD_DB.Where("id = ?", clusterID).First(clusterIns).Error; err != nil {
			global.DYCLOUD_LOG.Error("获取失败: " + err.Error())
			return nil, err
		}
		// Use reflection to check for specific types that need a different k8sType
		reqType := reflect.TypeOf(req).Elem()
		switch reqType {
		case reflect.TypeOf(nodes.NodeMetricsReq{}):
			k8sType = 2
		}
	}
	//else {
	//	return nil, errors.New("invalid request type")
	//}
	var sysUser = &system.SysUser{}
	if err := global.DYCLOUD_DB.Where("uuid = ?", uuid).First(sysUser).Error; err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}

	user, err := GetClusterByUserUUID(int(clusterIns.ID), uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("get user info failed: " + err.Error())
		return nil, err
	}
	isAdmin := false
	if sysUser.ID == 1 {
		isAdmin = true
	}
	var kubernetes *Kubernetes

	if k8sType == 1 {
		kubernetes = NewKubernetes(clusterIns, &user, isAdmin)
		d.Kubernetes = kubernetes
	} else {
		d.NodeMetrics = NewNodeMetrics(clusterIns, &user, isAdmin)
	}

	return kubernetes, nil
}

// GetClusterByUserUUID 根据用户UUID获取集群用户信息
//
// @Description 根据用户UUID获取集群用户信息
// @Tags Kubernetes
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param uuid path string true "用户UUID"
// @Success 200 {object} cluster2.User
// @Router /kubernetes/user/cluster [get]
func GetClusterByUserUUID(id int, uuid uuid.UUID) (user cluster2.User, err error) {
	err = global.DYCLOUD_DB.Where("cluster_id = ? and uuid = ?", id, uuid).First(&user).Error
	return
}
