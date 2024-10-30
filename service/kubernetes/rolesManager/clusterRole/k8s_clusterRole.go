package clusterRole

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/clusterrole"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sClusterRoleService struct {
	kubernetes.BaseService
}

func (k *K8sClusterRoleService) GetClusterRoleList(req clusterrole.GetClusterRoleListReq, uuid uuid.UUID) (*[]v1.ClusterRole, int, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	options := metav1.ListOptions{LabelSelector: req.LabelSelector}

	data, err := client.RbacV1().ClusterRoles().List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterClusterRole []v1.ClusterRole
	if req.Keyword != "" {
		for _, clusterRole := range data.Items {
			if strings.Contains(clusterRole.Name, req.Keyword) {
				filterClusterRole = append(filterClusterRole, clusterRole)
			}
		}
	} else {
		filterClusterRole = data.Items
	}

	result, err := paginate.Paginate(filterClusterRole, req.Page, req.PageSize)

	return result.(*[]v1.ClusterRole), len(filterClusterRole), nil
}
func (k *K8sClusterRoleService) DescribeClusterRole(req clusterrole.DescribeClusterRoleReq, uuid uuid.UUID) (*v1.ClusterRole, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	ClusterRole, err := client.RbacV1().ClusterRoles().Get(context.TODO(), req.ClusterRoleName, metav1.GetOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	return ClusterRole, nil
}
func (k *K8sClusterRoleService) UpdateClusterRole(req clusterrole.UpdateClusterRoleReq, uuid uuid.UUID) (*v1.ClusterRole, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	data, err := json.Marshal(req.Content)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	tmp := string(data)
	var ClusterRoleIns *v1.ClusterRole
	err = json.Unmarshal([]byte(tmp), &ClusterRoleIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.RbacV1().ClusterRoles().Update(context.TODO(), ClusterRoleIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sClusterRoleService) CreateClusterRole(req clusterrole.CreateClusterRoleReq, uuid uuid.UUID) (*v1.ClusterRole, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	data, err := json.Marshal(req.Content)
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	var ClusterRole *v1.ClusterRole
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &ClusterRole)
	ins, err := client.RbacV1().ClusterRoles().Create(context.TODO(), ClusterRole, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
func (k *K8sClusterRoleService) DeleteClusterRole(req clusterrole.DeleteClusterRoleReq, uuid uuid.UUID) error {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return err
	}
	err = client.RbacV1().ClusterRoles().Delete(context.TODO(), req.ClusterRoleName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
