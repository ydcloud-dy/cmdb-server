package role

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/roles"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sRoleService struct {
	kubernetes.BaseService
}

func (k *K8sRoleService) GetRoleList(req roles.GetRolesListReq, uuid uuid.UUID) (*[]v1.Role, int, error) {
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

	data, err := client.RbacV1().Roles(req.Namespace).List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterRole []v1.Role
	if req.Keyword != "" {
		for _, PV := range data.Items {
			if strings.Contains(PV.Name, req.Keyword) {
				filterRole = append(filterRole, PV)
			}
		}
	} else {
		filterRole = data.Items
	}

	result, err := paginate.Paginate(filterRole, req.Page, req.PageSize)

	return result.(*[]v1.Role), len(filterRole), nil
}
func (k *K8sRoleService) DescribeRole(req roles.DescribeRolesReq, uuid uuid.UUID) (*v1.Role, error) {
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
	Role, err := client.RbacV1().Roles(req.Namespace).Get(context.TODO(), req.RolesName, metav1.GetOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	return Role, nil
}
func (k *K8sRoleService) UpdateRole(req roles.UpdateRolesReq, uuid uuid.UUID) (*v1.Role, error) {
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
	var RoleIns *v1.Role
	err = json.Unmarshal([]byte(tmp), &RoleIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.RbacV1().Roles(req.Namespace).Update(context.TODO(), RoleIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sRoleService) CreateRole(req roles.CreateRolesReq, uuid uuid.UUID) (*v1.Role, error) {
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
	var Role *v1.Role
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &Role)
	ins, err := client.RbacV1().Roles(req.Namespace).Create(context.TODO(), Role, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
func (k *K8sRoleService) DeleteRole(req roles.DeleteRolesReq, uuid uuid.UUID) error {
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
	err = client.RbacV1().Roles(req.Namespace).Delete(context.TODO(), req.RolesName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
