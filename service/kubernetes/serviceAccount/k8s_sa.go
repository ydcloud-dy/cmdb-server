package serviceAccount

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/serviceAccount"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sServiceAccountService struct {
	kubernetes.BaseService
}

func (k *K8sServiceAccountService) GetServiceAccountList(req serviceAccount.GetServiceAccountReq,
	uuid uuid.UUID) (*[]v1.ServiceAccount, int, error) {
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
	data, err := client.CoreV1().ServiceAccounts(req.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var serviceAccountList []v1.ServiceAccount
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				serviceAccountList = append(serviceAccountList, item)
			}
		}
	} else {
		serviceAccountList = data.Items
	}

	result, err := paginate.Paginate(serviceAccountList, req.Page, req.PageSize)
	return result.(*[]v1.ServiceAccount), len(serviceAccountList), nil
}
func (k *K8sServiceAccountService) DescribeServiceAccount(req serviceAccount.DescribeServiceAccountReq, uuid uuid.UUID) (*v1.ServiceAccount, error) {
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
	ServiceAccount, err := client.CoreV1().ServiceAccounts(req.Namespace).Get(context.TODO(), req.ServiceAccountName, metav1.GetOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	return ServiceAccount, nil
}
func (k *K8sServiceAccountService) UpdateServiceAccount(req serviceAccount.UpdateServiceAccountReq, uuid uuid.UUID) (*v1.ServiceAccount, error) {
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
	var ServiceAccountIns *v1.ServiceAccount
	err = json.Unmarshal([]byte(tmp), &ServiceAccountIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.CoreV1().ServiceAccounts(req.Namespace).Update(context.TODO(), ServiceAccountIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sServiceAccountService) CreateServiceAccount(req serviceAccount.CreateServiceAccountReq, uuid uuid.UUID) (*v1.ServiceAccount, error) {
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
	var ServiceAccount *v1.ServiceAccount
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &ServiceAccount)
	ins, err := client.CoreV1().ServiceAccounts(req.Namespace).Create(context.TODO(), ServiceAccount, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
func (k *K8sServiceAccountService) DeleteServiceAccount(req serviceAccount.DeleteServiceAccountReq, uuid uuid.UUID) error {
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
	err = client.CoreV1().ServiceAccounts(req.Namespace).Delete(context.TODO(), req.ServiceAccountName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
