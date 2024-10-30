package clusterRoleBinding

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/clusterolebinding"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sClusterRoleBindingService struct {
	kubernetes.BaseService
}

func (k *K8sClusterRoleBindingService) GetClusterRoleBindingList(req clusterolebinding.GetClusterRoleBindingListReq, uuid uuid.UUID) (*[]v1.ClusterRoleBinding, int, error) {
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

	data, err := client.RbacV1().ClusterRoleBindings().List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterClusterRoleBinding []v1.ClusterRoleBinding
	if req.Keyword != "" {
		for _, PV := range data.Items {
			if strings.Contains(PV.Name, req.Keyword) {
				filterClusterRoleBinding = append(filterClusterRoleBinding, PV)
			}
		}
	} else {
		filterClusterRoleBinding = data.Items
	}

	result, err := paginate.Paginate(filterClusterRoleBinding, req.Page, req.PageSize)

	return result.(*[]v1.ClusterRoleBinding), len(filterClusterRoleBinding), nil
}
func (k *K8sClusterRoleBindingService) DescribeClusterRoleBinding(req clusterolebinding.DescribeClusterRoleBindingReq, uuid uuid.UUID) (*v1.ClusterRoleBinding, error) {
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
	ClusterRoleBinding, err := client.RbacV1().ClusterRoleBindings().Get(context.TODO(), req.ClusterRoleBindingName, metav1.GetOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	return ClusterRoleBinding, nil
}
func (k *K8sClusterRoleBindingService) UpdateClusterRoleBinding(req clusterolebinding.UpdateClusterRoleBindingReq, uuid uuid.UUID) (*v1.ClusterRoleBinding, error) {
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
	var ClusterRoleBindingIns *v1.ClusterRoleBinding
	err = json.Unmarshal([]byte(tmp), &ClusterRoleBindingIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.RbacV1().ClusterRoleBindings().Update(context.TODO(), ClusterRoleBindingIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sClusterRoleBindingService) CreateClusterRoleBinding(req clusterolebinding.CreateClusterRoleBindingReq, uuid uuid.UUID) (*v1.ClusterRoleBinding, error) {
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
	var ClusterRoleBinding *v1.ClusterRoleBinding
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &ClusterRoleBinding)
	ins, err := client.RbacV1().ClusterRoleBindings().Create(context.TODO(), ClusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
func (k *K8sClusterRoleBindingService) DeleteClusterRoleBinding(req clusterolebinding.DeleteClusterRoleBindingReq, uuid uuid.UUID) error {
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
	err = client.RbacV1().ClusterRoleBindings().Delete(context.TODO(), req.ClusterRoleBindingName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
