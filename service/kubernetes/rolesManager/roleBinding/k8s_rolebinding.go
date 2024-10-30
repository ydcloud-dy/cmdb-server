package roleBinding

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/roleBinding"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sRoleBindingService struct {
	kubernetes.BaseService
}

func (k *K8sRoleBindingService) GetRoleBindingList(req roleBinding.GetRoleBindingListReq, uuid uuid.UUID) (*[]v1.RoleBinding, int, error) {
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

	data, err := client.RbacV1().RoleBindings(req.Namespace).List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterRoleBinding []v1.RoleBinding
	if req.Keyword != "" {
		for _, PV := range data.Items {
			if strings.Contains(PV.Name, req.Keyword) {
				filterRoleBinding = append(filterRoleBinding, PV)
			}
		}
	} else {
		filterRoleBinding = data.Items
	}

	result, err := paginate.Paginate(filterRoleBinding, req.Page, req.PageSize)

	return result.(*[]v1.RoleBinding), len(filterRoleBinding), nil
}
func (k *K8sRoleBindingService) DescribeRoleBinding(req roleBinding.DescribeRoleBindingReq, uuid uuid.UUID) (*v1.RoleBinding, error) {
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
	RoleBinding, err := client.RbacV1().RoleBindings(req.Namespace).Get(context.TODO(), req.RoleBindingName, metav1.GetOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	return RoleBinding, nil
}
func (k *K8sRoleBindingService) UpdateRoleBinding(req roleBinding.UpdateRoleBindingReq, uuid uuid.UUID) (*v1.RoleBinding, error) {
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
	var RoleBindingIns *v1.RoleBinding
	err = json.Unmarshal([]byte(tmp), &RoleBindingIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.RbacV1().RoleBindings(req.Namespace).Update(context.TODO(), RoleBindingIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sRoleBindingService) CreateRoleBinding(req roleBinding.CreateRoleBindingReq, uuid uuid.UUID) (*v1.RoleBinding, error) {
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
	var RoleBinding *v1.RoleBinding
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &RoleBinding)
	ins, err := client.RbacV1().RoleBindings(req.Namespace).Create(context.TODO(), RoleBinding, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
func (k *K8sRoleBindingService) DeleteRoleBinding(req roleBinding.DeleteRoleBindingReq, uuid uuid.UUID) error {
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
	err = client.RbacV1().RoleBindings(req.Namespace).Delete(context.TODO(), req.RoleBindingName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
