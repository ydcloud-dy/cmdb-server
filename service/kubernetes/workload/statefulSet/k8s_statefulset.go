package statefulSet

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/statefulSet"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"github.com/goccy/go-json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sStatefulSetService struct {
	kubernetes.BaseService
}

func (k *K8sStatefulSetService) GetStatefulSetList(req statefulSet.GetStatefulSetListReq, uuid uuid.UUID) (*[]v1.StatefulSet, int, error) {
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

	data, err := client.AppsV1().StatefulSets(req.Namespace).List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterStatefulSet []v1.StatefulSet
	if req.Keyword != "" {
		for _, PV := range data.Items {
			if strings.Contains(PV.Name, req.Keyword) {
				filterStatefulSet = append(filterStatefulSet, PV)
			}
		}
	} else {
		filterStatefulSet = data.Items
	}

	result, err := paginate.Paginate(filterStatefulSet, req.Page, req.PageSize)

	return result.(*[]v1.StatefulSet), len(filterStatefulSet), nil
}
func (k *K8sStatefulSetService) DescribeStatefulSet(req statefulSet.DescribeStatefulSetReq, uuid uuid.UUID) (*v1.StatefulSet, error) {
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
	statefulset, err := client.AppsV1().StatefulSets(req.Namespace).Get(context.TODO(), req.StatefulSetName, metav1.GetOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	return statefulset, nil
}
func (k *K8sStatefulSetService) UpdateStatefulSet(req statefulSet.UpdateStatefulSetReq, uuid uuid.UUID) (*v1.StatefulSet, error) {
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
	var statefulSetIns *v1.StatefulSet
	err = json.Unmarshal([]byte(tmp), &statefulSetIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.AppsV1().StatefulSets(req.Namespace).Update(context.TODO(), statefulSetIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sStatefulSetService) CreateStatefulSet(req statefulSet.CreateStatefulSetReq, uuid uuid.UUID) (*v1.StatefulSet, error) {
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
	var statefulset *v1.StatefulSet
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &statefulset)
	ins, err := client.AppsV1().StatefulSets(req.Namespace).Create(context.TODO(), statefulset, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
func (k *K8sStatefulSetService) DeleteStatefulSet(req statefulSet.DeleteStatefulSetReq, uuid uuid.UUID) error {
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
	err = client.AppsV1().StatefulSets(req.Namespace).Delete(context.TODO(), req.StatefulSetName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
