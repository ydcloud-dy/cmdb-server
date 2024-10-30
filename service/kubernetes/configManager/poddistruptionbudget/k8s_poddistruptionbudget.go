package poddistruptionbudget

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/poddisruptionbudget"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sPoddisruptionbudgetService struct {
	kubernetes.BaseService
}

func (k *K8sPoddisruptionbudgetService) GetPoddisruptionbudgetList(req poddisruptionbudget.GetPoddisruptionbudgetListReq, uuid uuid.UUID) (*[]v1.PodDisruptionBudget, int, error) {
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
	data, err := client.PolicyV1().PodDisruptionBudgets(req.Namespace).List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterPoddisruptionbudget []v1.PodDisruptionBudget
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				filterPoddisruptionbudget = append(filterPoddisruptionbudget, item)
			}
		}
	} else {
		filterPoddisruptionbudget = data.Items
	}

	result, err := paginate.Paginate(filterPoddisruptionbudget, req.Page, req.PageSize)

	return result.(*[]v1.PodDisruptionBudget), len(filterPoddisruptionbudget), nil
}
func (k *K8sPoddisruptionbudgetService) DescribePoddisruptionbudget(req poddisruptionbudget.DescribePoddisruptionbudgetReq, uuid uuid.UUID) (*v1.PodDisruptionBudget, error) {
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
	PoddisruptionbudgetIns, err := client.PolicyV1().PodDisruptionBudgets(req.Namespace).Get(context.TODO(), req.PoddisruptionbudgetName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return PoddisruptionbudgetIns, nil
}
func (k *K8sPoddisruptionbudgetService) UpdatePoddisruptionbudget(req poddisruptionbudget.UpdatePoddisruptionbudgetReq, uuid uuid.UUID) (*v1.PodDisruptionBudget, error) {
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
	var PoddisruptionbudgetIns *v1.PodDisruptionBudget
	err = json.Unmarshal([]byte(tmp), &PoddisruptionbudgetIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.PolicyV1().PodDisruptionBudgets(req.Namespace).Update(context.TODO(), PoddisruptionbudgetIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sPoddisruptionbudgetService) DeletePoddisruptionbudget(req poddisruptionbudget.DeletePoddisruptionbudgetReq, uuid uuid.UUID) error {
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

	err = client.PolicyV1().PodDisruptionBudgets(req.Namespace).Delete(context.TODO(), req.PoddisruptionbudgetName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())

		return err
	}
	return nil
}
func (k *K8sPoddisruptionbudgetService) CreatePoddisruptionbudget(req poddisruptionbudget.CreatePoddisruptionbudgetReq, uuid uuid.UUID) (*v1.PodDisruptionBudget, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	data, err := json.Marshal(req.Content)
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	var Poddisruptionbudget *v1.PodDisruptionBudget
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &Poddisruptionbudget)
	ins, err := client.PolicyV1().PodDisruptionBudgets(req.Namespace).Create(context.TODO(), Poddisruptionbudget, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
