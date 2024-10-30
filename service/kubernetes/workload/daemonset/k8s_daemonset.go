package daemonset

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/daemonset"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sDaemonSetService struct {
	kubernetes.BaseService
}

func (k *K8sDaemonSetService) GetDaemonSetList(req daemonset.GetDaemonSetListReq, uuid uuid.UUID) (*[]v1.DaemonSet, int, error) {
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
	options := metav1.ListOptions{LabelSelector: req.LabelSelector, FieldSelector: req.FieldSelector}

	data, err := client.AppsV1().DaemonSets(req.Namespace).List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterDaemonsets []v1.DaemonSet
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				filterDaemonsets = append(filterDaemonsets, item)
			}
		}
	} else {
		filterDaemonsets = data.Items
	}

	result, err := paginate.Paginate(filterDaemonsets, req.Page, req.PageSize)

	return result.(*[]v1.DaemonSet), len(filterDaemonsets), nil
}
func (k *K8sDaemonSetService) DescribeDaemonSet(req daemonset.DescribeDaemonSetReq, uuid uuid.UUID) (*v1.DaemonSet, error) {
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
	daemonsetIns, err := client.AppsV1().DaemonSets(req.Namespace).Get(context.TODO(), req.DaemonsetName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return daemonsetIns, nil
}
func (k *K8sDaemonSetService) UpdateDaemonSet(req daemonset.UpdateDaemonSetReq, uuid uuid.UUID) (*v1.DaemonSet, error) {
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
	var daemonsetIns *v1.DaemonSet
	err = json.Unmarshal([]byte(tmp), &daemonsetIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.AppsV1().DaemonSets(req.Namespace).Update(context.TODO(), daemonsetIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sDaemonSetService) DeleteDaemonSet(req daemonset.DeleteDaemonSetReq, uuid uuid.UUID) error {
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

	err = client.AppsV1().DaemonSets(req.Namespace).Delete(context.TODO(), req.DaemonsetName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())

		return err
	}
	return nil
}
func (k *K8sDaemonSetService) CreateDaemonSet(req daemonset.CreateDaemonSetReq, uuid uuid.UUID) (*v1.DaemonSet, error) {
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
	var daemonset *v1.DaemonSet
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &daemonset)
	ins, err := client.AppsV1().DaemonSets(req.Namespace).Create(context.TODO(), daemonset, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
